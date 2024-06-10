package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/cache"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/utils"
	"go-protector/server/internal/utils/email"
	"sync"
	"time"
)

var loginPolicyServiceMap sync.Map
var policyServiceOnce sync.Once

type SysLoginPolicy struct {
	base.Service
}

var (
	defaultPolicySlice  []entity.SysLoginPolicy
	defaultGlobalPolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyGlobal,
		PolicyName: "全局认证策略",
		Enable:     "0",
	}
	defaultEmailPolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyEmail,
		PolicyName: "邮箱认证策略",
		Enable:     "0",
	}
	defaultOtpPolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyOtp,
		PolicyName: "动态密码认证策略",
		Enable:     "0",
	}
	defaultSharePolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyShare,
		PolicyName: "共享登录校验策略",
		Enable:     "0",
	}
	defaultIntruderPolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyIntruder,
		PolicyName: "防爆破登录策略",
		Enable:     "0",
	}
)

func init() {
	policyServiceOnce.Do(func() {
		loginPolicyServiceMap.Store(consts.LoginPolicyOtp, &otpLoginPolicyService{})
		loginPolicyServiceMap.Store(consts.LoginPolicyEmail, &emailLoginPolicyService{})
		defaultParam := map[string]interface{}{
			"mode": "1", // 0: 校验全部,1:只校验一个策略即可
		}
		bytes, _ := json.Marshal(defaultParam)
		defaultGlobalPolicy.Json = string(bytes)

		defaultParam = map[string]interface{}{
			"expireTime": 5, // 分钟
		}
		bytes, _ = json.Marshal(defaultParam)
		defaultEmailPolicy.Json = string(bytes)

		defaultParam = map[string]interface{}{
			"issuer":     "default-issuer", // 签发名称
			"period":     5,                // 秒
			"secretSize": 24,               // 秘钥长度
		}
		bytes, _ = json.Marshal(defaultParam)
		defaultOtpPolicy.Json = string(bytes)
		bytes, _ = json.Marshal(dto.ShareLoginPolicyDTO{
			LoginInterval:       0,
			SingleOnline:        0,
			SingleOnlineOperate: 0,
		})
		defaultSharePolicy.Json = string(bytes)

		bytes, _ = json.Marshal(dto.IntruderLoginPolicyDTO{
			FailThreshold: 5,
		})
		defaultIntruderPolicy.Json = string(bytes)
		defaultPolicySlice = append(defaultPolicySlice,
			defaultGlobalPolicy,
			defaultEmailPolicy,
			defaultOtpPolicy,
			defaultSharePolicy,
			defaultIntruderPolicy,
		)
	})

}

// Info 查询策略信息
func (_self *SysLoginPolicy) Info() (result *base.Result) {
	info := map[c_type.LoginPolicyCode]map[string]interface{}{}
	var modelSlice []entity.SysLoginPolicy
	if err := _self.DB.Table(table_name.SysLoginPolicy).Find(&modelSlice).Error; err != nil {
		return base.ResultFailureErr(err)
	}

	var model entity.SysLoginPolicy
	var valMap map[string]interface{}
	if len(modelSlice) == len(dto.LoginPolicyDTOFactory) {
		for _, model = range modelSlice {
			valMap = make(map[string]interface{})
			if err := json.Unmarshal([]byte(model.Json), &valMap); err != nil {
				return base.ResultFailureErr(err)
			}
			valMap["name"] = model.PolicyName
			valMap["enable"] = model.Enable
			info[model.PolicyCode] = valMap
			model = entity.SysLoginPolicy{}
		}
	} else {
		for _, model = range defaultPolicySlice {
			if err := _self.DB.Where(entity.SysLoginPolicy{PolicyCode: model.PolicyCode}).
				Attrs(model).FirstOrCreate(&model).Error; err != nil {
				return base.ResultFailureErr(err)
			}
			valMap = make(map[string]interface{})
			if err := json.Unmarshal([]byte(model.Json), &valMap); err != nil {
				return base.ResultFailureErr(err)
			}
			valMap["name"] = model.PolicyName
			valMap["enable"] = model.Enable
			info[model.PolicyCode] = valMap
			model = entity.SysLoginPolicy{}
		}
	}

	/*
		if err := _self.DB.Where(entity.SysLoginPolicy{PolicyCode: consts.LoginPolicyGlobal}).
			Assign(defaultGlobalPolicy).FirstOrCreate(&model).Error; err != nil {
			return base.ResultFailureErr(err)
		}
		if err := json.Unmarshal([]byte(model.PJson), &valMap); err != nil {
			return base.ResultFailureErr(err)
		}
		valMap["name"] = model.PolicyName
		info[consts.LoginPolicyGlobal] = valMap
		clear(valMap)
		model = entity.SysLoginPolicy{}

		if err := _self.DB.Where(entity.SysLoginPolicy{PolicyCode: consts.LoginPolicyEmail}).
			Assign(defaultEmailPolicy).FirstOrCreate(&model).Error; err != nil {
			return base.ResultFailureErr(err)
		}
		if err := json.Unmarshal([]byte(model.PJson), &valMap); err != nil {
			return base.ResultFailureErr(err)
		}
		valMap["name"] = model.PolicyName
		info[consts.LoginPolicyEmail] = valMap
		clear(valMap)
	*/
	return base.ResultSuccess(info)
}

// Save 保存
func (_self *SysLoginPolicy) Save(param map[c_type.LoginPolicyCode]map[string]interface{}) (
	result *base.Result) {

	if param == nil || len(param) <= 0 {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	// 只保存是否启用和json
	db := _self.DB.Begin()

	defer func() {
		if result.IsSuccess() {
			db.Commit()
			return
		}
		db.Rollback()
	}()

	for k, mapVal := range param {
		// 默认关闭
		enable := "0"
		if val, ok := mapVal["enable"]; ok {
			//interface 转 int
			enable = fmt.Sprintf("%v", val)
		}
		marshal, err := json.Marshal(mapVal)
		if err != nil {
			return base.ResultFailureErr(err)
		}

		updateMap := map[string]interface{}{
			"enable": enable,
			"json":   string(marshal),
		}
		policyDTO, err := dto.NewLoginPolicyDTO(k, mapVal)
		if err = policyDTO.Validate(policyDTO); err != nil {
			return base.ResultFailureErr(err)
		}
		if err = db.Table(table_name.SysLoginPolicy).
			Where("policy_code = ?", k).Updates(updateMap).Error; err != nil {
			return base.ResultFailureErr(err)
		}
	}
	return base.ResultSuccessMsg("保存成功")

}

// ProcessLoginPolicy 处理登录策略
func (_self *SysLoginPolicy) ProcessLoginPolicy(policyParam *dto.LoginPolicyParamDTO, user *entity.SysUser,
	policyDTOMap PolicyDTOMap) (res *base.Result) {
	if policyParam == nil || user == nil {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	var policyDTO dto.ILoginPolicyDTO
	var err error
	var ok bool
	if policyDTO, ok = policyDTOMap[policyParam.PolicyCode]; !ok {
		return base.ResultFailureMsg("获取策略失败")
	}
	if err != nil {
		return base.ResultFailureErr(err)
	}
	var mailService ILoginPolicyOfMailService
	mailService, err = GetLoginPolicyOfMailService(policyParam.PolicyCode)
	if err != nil {
		return base.ResultFailureErr(err)
	}

	processDTO := &dto.LoginPolicyProcessDTO{
		Service:   _self.Service,
		PolicyDTO: policyDTO,
		ParamDTO:  policyParam,
		User: &current.User{
			ID:        user.ID,
			SessionId: user.SessionId,
			LoginName: user.LoginName,
			Email:     user.Email,
		},
	}
	if policyParam.Operate == 0 {
		if err = mailService.Send(processDTO); err != nil {
			res = base.ResultFailureErr(err)
		} else {
			res = base.ResultSuccessMsg("发送成功")
		}
	} else {
		if err = mailService.Verify(processDTO); err != nil {
			res = base.ResultFailureErr(err)
		} else {
			res = base.ResultSuccessMsg("校验成功")
		}
	}
	return
}

// GetLoginPolicyOfMailService 获取与邮箱有关的登录策略
func GetLoginPolicyOfMailService(policyCode c_type.LoginPolicyCode) (
	policyService ILoginPolicyOfMailService, err error) {

	value, ok := loginPolicyServiceMap.Load(policyCode)
	if !ok {
		err = errors.New("无法获取策略")
	}
	if policyService, ok = value.(ILoginPolicyOfMailService); !ok {
		err = errors.New("无法获取策略")
	}
	return
}

type ILoginPolicyOfMailService interface {
	Send(processDTO *dto.LoginPolicyProcessDTO) error
	// Verify email,otp
	// isPass 0:未通过,1:通过, 只有校验的时候才使用 isPass字段
	// err 错误信息
	Verify(processDTO *dto.LoginPolicyProcessDTO) error
}

type emailLoginPolicyService struct {
}

func (_self emailLoginPolicyService) Send(processDTO *dto.LoginPolicyProcessDTO) (err error) {
	//发送邮箱验证码 防止暴力发送
	if processDTO == nil {
		return c_error.ErrParamInvalid
	}
	if processDTO.User.Email == "" {
		err = errors.New("用户未配置邮箱,请联系管理员")
		return
	}
	// 1分钟
	redisClient := cache.GetRedisClient()
	keyPre := fmt.Sprintf(consts.LoginPolicyCacheKeyFmt, processDTO.User.LoginName, processDTO.ParamDTO.SessionId)
	redisKey := fmt.Sprintf("%s:%s", keyPre, "email")
	lockKey := fmt.Sprintf("%s:%s", redisKey, "lock")

	nx := redisClient.SetNX(processDTO.Service.Context, lockKey, "lock", time.Minute)
	var isSuccess bool
	isSuccess, err = nx.Result()
	if err != nil {
		return
	}
	if !isSuccess {
		err = errors.New("频繁发送过于频繁,请稍后再试")
		return
	}
	defer redisClient.Del(processDTO.Service.Context, lockKey)

	emailDTO, ok := processDTO.PolicyDTO.(*dto.EmailLoginPolicyDTO)
	if !ok {
		return c_error.ErrParamInvalid
	}

	expireTime := emailDTO.ExpireTime
	expiration := time.Duration(expireTime) * time.Minute
	if err = email.VerifySendInterval(processDTO.Service.Context, redisKey, expiration, 1*time.Minute); err != nil {
		return err
	}

	randomCode := utils.GetRandomCode()

	body := fmt.Sprintf("您好: 您的验证码为: %s,有效期为%d分钟", randomCode, expireTime)

	if err = email.Send(&email.SendDTO{
		To:      processDTO.User.Email,
		Subject: "登录验证码",
		Body:    body,
	}); err != nil {
		processDTO.Service.Logger.Error("send: {},err: %s", processDTO.User.Email, err)
		return
	}

	err = redisClient.Set(processDTO.Service.Context,
		redisKey, randomCode, expiration).Err()

	return
}

func (_self emailLoginPolicyService) Verify(processDTO *dto.LoginPolicyProcessDTO) (err error) {
	redisClient := cache.GetRedisClient()
	keyPre := fmt.Sprintf(consts.LoginPolicyCacheKeyFmt, processDTO.User.LoginName, processDTO.ParamDTO.SessionId)
	redisKey := fmt.Sprintf("%s:%s", keyPre, "email")
	var val string
	val, err = redisClient.Get(processDTO.Service.Context, redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = errors.New("验证码无效或已过期")
		}
		return
	}
	if val != processDTO.ParamDTO.Val {
		err = errors.New("验证码无效或已过期")
	}
	return

}

type otpLoginPolicyService struct {
}

func (_self otpLoginPolicyService) Send(processDTO *dto.LoginPolicyProcessDTO) (err error) {
	//发送邮箱验证码 防止暴力发送
	if processDTO == nil {
		return c_error.ErrParamInvalid
	}
	if processDTO.User.Email == "" {
		err = errors.New("用户未配置邮箱,请联系管理员")
		return
	}
	otpDTO, ok := processDTO.PolicyDTO.(*dto.OTPLoginPolicyDTO)
	if !ok {
		return c_error.ErrParamInvalid
	}
	// 1分钟
	redisClient := cache.GetRedisClient()
	keyPre := fmt.Sprintf(consts.LoginPolicyCacheKeyFmt, processDTO.User.LoginName, processDTO.ParamDTO.SessionId)
	redisKey := fmt.Sprintf("%s:%s", keyPre, "otp")
	lockKey := fmt.Sprintf("%s:%s", redisKey, "lock")

	nx := redisClient.SetNX(processDTO.Service.Context, lockKey, "lock", time.Minute)
	var isSuccess bool
	isSuccess, err = nx.Result()
	if err != nil {
		return
	}
	if !isSuccess {
		err = errors.New("频繁发送过于频繁,请稍后再试")
		return
	}
	defer redisClient.Del(processDTO.Service.Context, lockKey)

	if err = email.VerifySendInterval(processDTO.Service.Context, redisKey, 1*time.Minute, 0); err != nil {
		return err
	}

	var otpService SysOtpBind
	processDTO.Service.MakeService(&otpService)
	var qrCodeBase64 string
	if qrCodeBase64, err = otpService.GetQRCodeBase64ByUser(*processDTO.User, *otpDTO); err != nil {
		return err
	}
	// http://www.cdiso.cn/article/ejppi.html
	err = email.SendImage(email.SendDTO{
		To:      processDTO.User.Email,
		Subject: "OTP认证信息",
		Body:    "请使用OTP认证器扫描一下二维码",
	}, qrCodeBase64)

	return

}

func (_self otpLoginPolicyService) Verify(processDTO *dto.LoginPolicyProcessDTO) (err error) {
	var otpService SysOtpBind
	processDTO.Service.MakeService(&otpService)
	err = otpService.VerifyCode(*processDTO.User, processDTO.ParamDTO.Val)
	return
}
