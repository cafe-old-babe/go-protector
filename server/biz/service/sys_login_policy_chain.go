package service

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dao"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/cache"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
	"sync"
	"time"
)

var once sync.Once

type PolicyDTOMap map[c_type.LoginPolicyCode]dto.ILoginPolicyDTO

type LoginPolicyHandler func(_self base.Service, sysUser *entity.SysUser, loginDTO *dto.LoginDTO,
	policyDTOMap PolicyDTOMap) (res *base.Result, isLoop bool)

var chain []LoginPolicyHandler

var (
	// checkUserStatus 校验用户状态与有效期
	checkUserStatus = func(_self base.Service, sysUser *entity.SysUser,
		loginDTO *dto.LoginDTO, policyDTOMap PolicyDTOMap) (res *base.Result, isLoop bool) {
		_self.GetLogger().Debug("chain do checkUserExpirationAt")
		if sysUser.UserStatus != 0 {
			res = base.ResultFailureErr(c_error.ErrLoginNameOrPasswordIncorrect)
		}
		// 4-23	【实战】保存用户信息-掌握Gin依赖校验、GO语言defer特性
		// 检查有效期
		if !sysUser.ExpirationAt.Valid {
			return
		}
		// 当前时间是否在有效期之后
		if !time.Now().After(sysUser.ExpirationAt.Time) {
			isLoop = true
			return
		}

		_self.GetLogger().Error("用户: %s 已过有效期", loginDTO.LoginName)
		res = base.ResultFailureMsg(c_error.ErrLoginNameOrPasswordIncorrect.Error())
		// 更新用户信息
		sysUser.UserStatus = consts.LockTypeExpire
		sysUser.LockReason = sql.NullString{
			String: "用户已过有效期",
			Valid:  true,
		}
		sysUser.UpdatedBy = sysUser.ID
		if err := dao.SysUser.LockUser(_self.GetDB(), sysUser); err != nil {
			_self.GetLogger().Error("用户: %s lockUser err: %v", loginDTO.LoginName, err)
		}
		return
	}

	// checkUserPassword 校验用户密码
	checkUserPassword = func(_self base.Service, sysUser *entity.SysUser,
		loginDTO *dto.LoginDTO, policyDTOMap PolicyDTOMap) (res *base.Result, isLoop bool) {
		_self.GetLogger().Debug("chain do checkUserPassword")
		if loginDTO.PolicyParam != nil {
			// 策略登录不参与校验密码
			return
		}
		if sysUser.Password == loginDTO.Password {
			return
		}
		res = base.ResultFailureMsg(c_error.ErrLoginNameOrPasswordIncorrect.Error())

		policyDTO, ok := policyDTOMap[consts.LoginPolicyIntruder]
		if !ok {
			_self.GetLogger().Debug("获取防爆破策略失败")
			return
		}
		if !policyDTO.IsEnable() {
			return
		}

		intruderLoginPolicyDTO := policyDTO.(*dto.IntruderLoginPolicyDTO)
		threshold := intruderLoginPolicyDTO.FailThreshold
		redisClient := cache.GetRedisClient()

		redisKey := fmt.Sprintf(consts.LoginIntruderCacheKeyFmt, time.Now().Day(), loginDTO.LoginName)
		redisClient.SetNX(_self.GetContext(), redisKey, 0, time.Hour*24)

		val, err := redisClient.Incr(_self.GetContext(), redisKey).Result()
		if err != nil {
			_self.GetLogger().Error("incr %s, err: %v", redisKey, err)
			return
		}
		failCount := uint(val)
		if failCount > threshold {
			// 锁定用户
			// 更新用户信息
			sysUser.UserStatus = consts.LockTypePasswordFailure
			newTime := c_type.NewTime(time.Now())
			sysUser.LockReason = sql.NullString{
				String: fmt.Sprintf("[%s]用户输入密码错误次数达到阈值[%d],系统锁定", newTime.String(), threshold),
				Valid:  true,
			}
			sysUser.UpdatedBy = sysUser.ID
			if err := dao.SysUser.LockUser(_self.GetDB(), sysUser); err != nil {
				_self.GetLogger().Error("用户: %s lockUser err: %v", loginDTO.LoginName, err)
			}
		}
		return

	}

	// checkLoginInterval 校验登录间隔
	checkLoginInterval = func(_self base.Service, sysUser *entity.SysUser,
		loginDTO *dto.LoginDTO, policyDTOMap PolicyDTOMap) (res *base.Result, isLoop bool) {
		policyDTO := policyDTOMap[consts.LoginPolicyShare]
		if !policyDTO.IsEnable() {
			isLoop = false
			return
		}
		isLoop = true
		var err error
		var ginCtx *gin.Context
		if ginCtx, err = _self.GetGinCtx(); err != nil {
			res = base.ResultFailureErr(err)
		}
		if sysUser.LastLoginIp == ginCtx.ClientIP() {
			return
		}
		if !sysUser.LastLoginTime.Valid {
			return
		}

		if sysUser, err = dao.SysUser.FindUserInfoByDTO(_self.GetDB(), &dto.FindUserDTO{
			ID: sysUser.ID,
		}); err != nil {
			res = base.ResultFailureErr(err)
			return
		}
		loginPolicyDTO := policyDTO.(*dto.ShareLoginPolicyDTO)
		lastLoginTime := sysUser.LastLoginTime.Time
		if !lastLoginTime.Add(time.Duration(loginPolicyDTO.LoginInterval)).After(time.Now()) {
			return
		}
		return
	}
	// 校验单用户登录策略
	checkLoginSingle = func(_self base.Service, sysUser *entity.SysUser,
		loginDTO *dto.LoginDTO, policyDTOMap PolicyDTOMap) (res *base.Result, isLoop bool) {

		policyDTO := policyDTOMap[consts.LoginPolicyShare]
		if !policyDTO.IsEnable() {
			isLoop = false
			return
		}
		loginPolicyDTO := policyDTO.(*dto.ShareLoginPolicyDTO)
		if loginPolicyDTO.SingleOnline == 0 {
			return
		}
		isLoop = true
		redisClient := cache.GetRedisClient()
		key := fmt.Sprintf(consts.OnlineUserCacheKeyFmt, loginDTO.LoginName, "*")
		scan := redisClient.Scan(_self.GetContext(), 0, key, 0)
		if scan.Err() != nil {
			res = base.ResultFailureErr(scan.Err())
			return
		}
		var keys []string
		iterator := scan.Iterator()
		for iterator.Next(_self.GetContext()) {
			if loginPolicyDTO.SingleOnlineOperate == 0 {
				res = base.ResultFailureMsg("该账号正在使用中,请稍后再试")
				return
			}
			keys = append(keys, iterator.Val())
		}

		if len(keys) <= 0 {
			return
		}
		value := _self.GetContext().Value("status")
		val, ok := value.(string)
		if !ok {
			return
		}
		if val == "second" {
			if _, err := redisClient.Del(_self.GetContext(), keys...).Result(); err != nil {
				res = base.ResultFailureMsg("系统繁忙,请稍后再试")
				return
			}
		}

		return
	}

	// checkAuthenticationPolicy 认证策略
	// 400:失败提示信息; 200:登录成功; 201:触发登录策略; 203:展示提示信息
	checkAuthenticationPolicy = func(_self base.Service, sysUser *entity.SysUser,
		loginDTO *dto.LoginDTO, policyDTOMap PolicyDTOMap) (res *base.Result, isLoop bool) {
		policyParam := loginDTO.PolicyParam
		var sysLoginPolicyService SysLoginPolicy
		_self.MakeService(&sysLoginPolicyService)
		res = sysLoginPolicyService.Info()
		if !res.IsSuccess() {
			return
		}

		var userService SysUser
		_self.MakeService(&userService)
		if policyParam == nil {
			res = userService.GetLoginPolicyResult(sysUser, policyDTOMap)
			return
		}
		res = userService.ValidateLoginPolicyParam(loginDTO, sysUser, policyDTOMap)
		return
	}
)

func init() {
	once.Do(func() {
		chain = append(chain, checkUserStatus)
		chain = append(chain, checkUserPassword)
		chain = append(chain, checkLoginInterval)
		chain = append(chain, checkLoginSingle)
		chain = append(chain, checkAuthenticationPolicy)
	})

}

type SysLoginPolicyChain struct {
	base.Service
}

func (_self *SysLoginPolicyChain) Do(sysUser *entity.SysUser, loginDTO *dto.LoginDTO, policyDTOMap PolicyDTOMap,
	loginSuccessFunc func(entity *entity.SysUser) (res *base.Result)) (res *base.Result) {

	var loopChain []LoginPolicyHandler
	var isLoop bool
	_self.Set("status", "first")

	for i := range chain {
		res, isLoop = chain[i](_self.Service, sysUser, loginDTO, policyDTOMap)
		if isLoop {
			loopChain = append(loopChain, chain[i])
		}
		if res == nil {
			continue
		}
		if !res.IsSuccess() {
			return
		}
	}
	for i := range loopChain {
		_self.Set("status", "second")
		res, _ = loopChain[i](_self.Service, sysUser, loginDTO, policyDTOMap)
		if res == nil || !res.IsSuccess() {
			continue
		}
	}
	return loginSuccessFunc(sysUser)
}
