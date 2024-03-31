package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-protector/server/core/base"
	"go-protector/server/core/cache"
	"go-protector/server/core/consts"
	"go-protector/server/core/consts/table_name"
	"go-protector/server/core/current"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_jwt"
	"go-protector/server/core/custom/c_type"
	"go-protector/server/core/scope"
	"go-protector/server/core/utils"
	"go-protector/server/dao"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"go-protector/server/models/vo"
	"gorm.io/gorm"
	"time"
)

type SysUser struct {
	base.Service
}

// FindUserByDTO 查询用户信息,如果没查到 返回失败结果
func (_self *SysUser) FindUserByDTO(findDTO *dto.FindUser) (res *base.Result) {
	var sysUser *entity.SysUser
	var err error
	sysUser, err = dao.SysUser.FindUserByDTO(_self.DB, findDTO)
	if err != nil {
		_self.Logger.Error("FindUserByDTO err: %v", err)
		//if errors.Is(c_error.ErrRecordNotFoundSysUser, err) {
		//	res = base.ResultFailureMsg(c_error.ErrRecordNotFoundSysUser.Error())
		//}
		res = base.ResultFailureMsg(err.Error())
		return
	}
	if sysUser == nil {
		res = base.ResultFailureMsg(c_error.ErrLoginNameOrPasswordIncorrect.Error())
		return
	}
	res = base.ResultSuccess(sysUser)
	return
}

// LoginSuccess 登录成功
func (_self *SysUser) LoginSuccess(entity *entity.SysUser) (res *base.Result) {
	var err error
	// 更新最后登录时间 最后登录IP
	if err = dao.SysUser.UpdateLastLoginIp(_self.DB, entity.ID, _self.Context.ClientIP()); err != nil {
		_self.Logger.Error("用户: %s UpdateLastLoginIp err: %v", entity.LoginName, err)
	}

	userDTO := &current.User{
		ID:        entity.ID,
		SessionId: entity.SessionId,
		LoginName: entity.LoginName,
		Email:     entity.Email,
		UserName:  entity.Username,
		LoginTime: time.Now().Format(time.DateTime),
		LoginIp:   _self.Context.ClientIP(),
		Avatar:    "https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png",
		DeptId:    entity.DeptId,
		RoleIds:   entity.SysRoleIds,
	}
	for _, role := range entity.SysRoles {
		if role.RoleType == 0 {
			userDTO.IsAdmin = true
			break
		}
	}
	// 生成Token
	jwtString, expireAt, err := c_jwt.GenerateToken(userDTO)

	res = base.ResultSuccess(dto.LoginSuccess{
		SysUser:  userDTO,
		Token:    *jwtString,
		ExpireAt: expireAt,
	})
	return
}

func (_self *SysUser) SetStatus(dto *dto.SetStatus) (err error) {
	if dto == nil || dto.ID <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	if dto.UserStatus <= 0 {
		// 解锁
		if err = dao.SysUser.UnlockUser(_self.DB, dto); err != nil {
			_self.Logger.Error("SetStatus UnlockUser err: %v", err)
		}
	} else {
		// 加锁
		if err = dao.SysUser.LockUser(_self.DB, &entity.SysUser{
			ModelId:    entity.ModelId{ID: dto.ID},
			LockReason: sql.NullString{Valid: len(dto.LockReason) > 0, String: dto.LockReason},
			UserStatus: dto.UserStatus,
		}); err != nil {
			_self.Logger.Error("SetStatus LockUser err: %v", err)
		}
	}

	return
}

// UserInfo https://pro.antdv.com/docs/authority-management
func (_self *SysUser) UserInfo() (res *base.Result) {
	// 查询用户角色
	user, ok := current.GetUser(_self.Context)
	if !ok {
		res = base.ResultFailureMsg("获取当前用户失败")
		return
	}
	// 查询角色关联菜单
	roleIds := user.RoleIds
	if len(roleIds) <= 0 {
		res = base.ResultFailureMsg("当前用户未绑定角色")
		return
	}
	var roleService SysRole
	_self.MakeService(&roleService)
	menuSlice, buttonSlice, err := roleService.GetMenuByRoleIds(roleIds, user.IsAdmin)
	if err != nil {
		return base.ResultFailureErr(err)
	}
	buttonMap := map[uint64][]entity.SysMenu{}
	for _, button := range buttonSlice {
		buttonMap[button.PID] = append(buttonMap[button.PID], button)
	}

	// 封装
	permissionSlice := make([]dto.Permission, 0)

	for _, menu := range menuSlice {

		permission := dto.Permission{
			PermissionId:    menu.Permission,
			PermissionName:  menu.Name,
			ActionEntitySet: make([]dto.ActionEntity, 0),
		}
		buttonSlice = buttonMap[menu.ID]
		for _, button := range buttonSlice {
			permission.ActionEntitySet =
				append(permission.ActionEntitySet,
					dto.ActionEntity{
						Action:       button.Permission,
						Describe:     button.Name,
						DefaultCheck: false,
					})
		}
		if actByte, err := json.Marshal(permission.ActionEntitySet); err == nil {
			permission.Actions = string(actByte)
		}

		permissionSlice = append(permissionSlice, permission)

	}
	// 对接antdPro
	roleInfo := dto.RoleInfo{
		Name: user.UserName,
		Role: dto.Role{
			Permissions: permissionSlice,
		},
	}
	return base.ResultSuccess(roleInfo)

}

// Nav 获取菜单
func (_self *SysUser) Nav() (res *base.Result) {
	// 查询用户角色
	user, ok := current.GetUser(_self.Context)
	if !ok {
		res = base.ResultFailureMsg("获取当前用户失败")
		return
	}
	// 查询角色关联菜单
	roleIds := user.RoleIds
	if len(roleIds) <= 0 {
		res = base.ResultFailureMsg("当前用户未绑定角色")
		return
	}

	menuSlice, err := dao.SysRole.GetPermissionSliceByIds(_self.DB, roleIds, []int8{0, 1}, user.IsAdmin)
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	var menuInfoSlice []dto.MenuInfo
	for _, menu := range menuSlice {

		menuInfoSlice = append(menuInfoSlice, dto.MenuInfo{
			Id:        menu.ID,
			ParentId:  menu.PID,
			Name:      menu.Permission,
			Path:      menu.Path,
			Component: menu.Component,
			Redirect:  nil,
			Meta: dto.MetaInfo{
				Title: menu.Name,
				Show:  menu.Hidden != 1,
			},
		})
	}
	return base.ResultSuccess(menuInfoSlice)
}

// Page 人员分页查询
func (_self *SysUser) Page(req *dto.UserPageReq) (result *base.Result) {
	if req == nil {
		result = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var count int64
	var voSlice []vo.UserPage
	tx := _self.DB.Table(table_name.SysUser + " as u")
	err := tx.Select([]string{
		"u.id",
		"u.login_name",
		"u.username",
		"u.user_status as status",
		"u.email",
		"u.sex",
		"d.id as dept_id",
		"d.dept_name",
		"p.post_names",
		"p.post_ids",
		"r.role_names",
		"r.role_ids",
	}).Scopes(
		scope.Paginate(req.GetPagination()),
		scope.Like("u.login_name", req.LoginName),
		scope.Like("u.username", req.Username),
		func(db *gorm.DB) *gorm.DB {
			if len(req.DeptIds) > 0 {
				db = db.Where("d.id in (?)", req.DeptIds)
			}
			return db
		},
	).Joins(`left join `+table_name.SysDept+` d on d.id = u.dept_id`).
		Joins("left join (?) as p on p.user_id = u.id", dao.SysPost.JoinUserPostDB(_self.DB)).
		Joins("left join (?) as r on r.user_id = u.id", dao.SysRole.JoinUserRoleDB(_self.DB)).
		Find(&voSlice).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		result = base.ResultFailureErr(err)
	} else {
		result = base.ResultPage(voSlice, req.GetPagination(), count)
	}
	return
}

func (_self *SysUser) Save(req *dto.UserSaveReq) (result *base.Result) {

	if err := dao.SysUser.Save(_self.DB, req); err != nil {
		return base.ResultFailureErr(err)
	}

	return base.ResultSuccessMsg("保存成功")

}

func (_self *SysUser) DeleteByIds(req *base.IdsReq) (result *base.Result) {

	if req == nil || len(req.Ids) <= 0 {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	db := _self.DB.Begin()
	defer func() {
		if result == nil || !result.IsSuccess() {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	// delete user
	if err := db.Delete(&entity.SysUser{}, req.Ids).Error; err != nil {
		result = base.ResultFailureErr(err)
		return
	}

	// delete post_relation
	if err := db.Delete(&entity.SysPostRelation{},
		"relation_id in (?) and relation_type = ?", req.Ids, consts.User).Error; err != nil {
		result = base.ResultFailureErr(err)
		return
	}

	// delete role_relation
	if err := db.Delete(&entity.SysRoleRelation{},
		"relation_id in (?) and relation_type = ?", req.Ids, consts.User).Error; err != nil {
		result = base.ResultFailureErr(err)
		return
	}

	return base.ResultSuccessMsg("删除成功")
}

// checkLogin 400:失败提示信息; 200:登录成功; 201:触发登录策略; 203:展示提示信息
// loginDTO 必传
// sysUser 可选
func (_self *SysUser) checkLogin(loginDTO *dto.LoginDTO, sysUser *entity.SysUser) (result *base.Result) {

	if loginDTO == nil || nil == sysUser {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}

	var sysLoginPolicyService SysLoginPolicy
	_self.MakeService(&sysLoginPolicyService)
	result = sysLoginPolicyService.Info()
	if !result.IsSuccess() {
		return
	}
	policyInfoMap, ok := result.Data.(map[c_type.LoginPolicyCode]map[string]interface{})
	if !ok {
		result = base.ResultFailureErr(errors.New("查询策略失败,请联系管理员"))
		return
	}
	var policyDTO dto.ILoginPolicyDTO
	var err error

	policyDTOMap := make(PolicyDTOMap)
	for k, v := range policyInfoMap {
		if policyDTO, err = dto.NewLoginPolicyDTO(k, v); err != nil {
			result = base.ResultFailureErr(err)
			return
		}
		policyDTOMap[k] = policyDTO
	}

	var chainService SysLoginPolicyChain
	_self.MakeService(&chainService)
	return chainService.Do(sysUser, loginDTO, policyDTOMap, _self.LoginSuccess)
	/*
	   	var err error
	   	now := time.Now()
	   	if sysUser.UserStatus != 0 {
	   		if sysUser.UserStatus != consts.LockTypePasswordFailure {
	   			result = base.ResultFailureErr(c_error.ErrLoginNameOrPasswordIncorrect)
	   			return
	   		}
	   	}

	   	// 检查有效期
	   	if sysUser.ExpirationAt.Valid {
	   		if now.After(sysUser.ExpirationAt.Time) {
	   			_self.Logger.Error("用户: %s 已过有效期", loginDTO.LoginName)
	   			result = base.ResultFailureMsg(c_error.ErrLoginNameOrPasswordIncorrect.Error())
	   			// 更新用户信息
	   			sysUser.UserStatus = consts.LockTypeExpire
	   			sysUser.LockReason = sql.NullString{
	   				String: "用户已过有效期",
	   				Valid:  true,
	   			}
	   			sysUser.UpdatedBy = sysUser.ID
	   			if err := dao.SysUser.LockUser(_self.DB, sysUser); err != nil {
	   				_self.Logger.Error("用户: %s lockUser err: %v", loginDTO.LoginName, err)
	   			}
	   			return
	   		}
	   	}
	   	// 校验密码
	   	if sysUser.Password != loginDTO.Password {
	   		result = base.ResultFailureMsg(c_error.ErrLoginNameOrPasswordIncorrect.Error())
	   		return
	   	}

	   	// 校验全局策略
	   	policyParam := loginDTO.PolicyParam
	   	var sysLoginPolicyService SysLoginPolicy
	   	_self.MakeService(&sysLoginPolicyService)
	   	result = sysLoginPolicyService.Info()
	   	if !result.IsSuccess() {
	   		return
	   	}
	   	policyInfoMap, ok := result.Data.(map[c_type.LoginPolicyCode]map[string]interface{})
	   	if !ok {
	   		result = base.ResultFailureErr(errors.New("查询策略失败,请联系管理员"))
	   		return
	   	}
	   	var policyDTO dto.ILoginPolicyDTO

	   	var policyDTOMap PolicyDTOMap
	   	for k, v := range policyInfoMap {
	   		if policyDTO, err = dto.NewLoginPolicyDTO(k, v); err != nil {
	   			result = base.ResultFailureErr(err)
	   			return
	   		}
	   		policyDTOMap[k] = policyDTO
	   		if k == consts.LoginPolicyGlobal {
	   			if !policyDTO.IsEnable() {
	   				goto doLoginSuccess
	   			}

	   		}

	   	}

	   	// 创建责任链

	   	if policyParam == nil {
	   		policyDTO, err = dto.NewLoginPolicyDTO(consts.LoginPolicyGlobal, policyInfoMap[consts.LoginPolicyGlobal])
	   		if err != nil {
	   			result = base.ResultFailureErr(err)
	   			return
	   		}
	   		if !policyDTO.IsEnable() {
	   			goto doLoginSuccess
	   		}
	   		result = _self.GetLoginPolicyResult(sysUser, policyDTOMap)
	   		if nil == result {
	   			goto doLoginSuccess
	   		}
	   		return
	   	}
	   	result = _self.ValidateLoginPolicyParam(loginDTO, sysUser, policyDTOMap)
	   	if result != nil {
	   		return
	   	}

	   	// 校验共享登录策略
	   	if policyDTO, err = dto.NewLoginPolicyDTO(consts.LoginPolicyShare, policyInfoMap[consts.LoginPolicyShare]); err != nil {
	   		result = base.ResultFailureErr(err)
	   		return
	   	}
	   	if !policyDTO.IsEnable() {
	   		goto doLoginSuccess
	   	}

	   doLoginSuccess:
	   	return _self.LoginSuccess(sysUser)*/
}

// GetLoginPolicyResult 获取全局认证策略结果
func (_self *SysUser) GetLoginPolicyResult(user *entity.SysUser, policyInfoMap PolicyDTOMap) (result *base.Result) {
	if nil == user {
		result = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var policyDTO dto.ILoginPolicyDTO
	var ok bool
	if policyDTO, ok = policyInfoMap[consts.LoginPolicyGlobal]; !ok {
		result = base.ResultFailureMsg("获取策略失败")
		return
	}
	if !policyDTO.IsEnable() {
		return nil
	}
	res := dto.LoginPolicyResultDTO{
		SessionId: utils.GetNextIdStr(),
		Mode:      fmt.Sprintf("%v", policyDTO.GetKey("mode")),
	}
	user.SessionId = res.SessionId
	for k, v := range policyInfoMap {
		if k == consts.LoginPolicyGlobal ||
			k == consts.LoginPolicyShare ||
			k == consts.LoginPolicyIntruder {
			continue
		}
		if v.IsEnable() {
			res.PolicyCode = append(res.PolicyCode, k)
		}
	}
	// 没有开启的策略
	if len(res.PolicyCode) <= 0 {
		return nil
	}

	// 录策略信息
	policyKey := fmt.Sprintf(consts.LoginPolicyCacheKeyFmt, user.LoginName, res.SessionId)
	redisClient := cache.GetRedisClient()
	var err error
	defer func() {
		if err != nil {
			redisClient.Del(_self.Context, policyKey)
			if result == nil {
				result = base.ResultFailureErr(err)
			}
		}
	}()
	var bytes []byte
	bytes, err = json.Marshal(res)
	if err != nil {
		result = base.ResultFailureErr(err)
		return
	}
	if err = redisClient.Set(_self.Context, policyKey, string(bytes), time.Minute*30).Err(); err != nil {
		return
	}

	if bytes, err = json.Marshal(user); err != nil {
		return
	}

	return base.ResultCustom(201, res, "请继续认证")
}

// ValidateLoginPolicyParam 校验认证策略
func (_self *SysUser) ValidateLoginPolicyParam(loginDTO *dto.LoginDTO, user *entity.SysUser, policyDTOMap PolicyDTOMap) (res *base.Result) {

	policyParam := loginDTO.PolicyParam
	if nil == policyParam || len(policyParam.SessionId) <= 0 {
		return base.ResultFailureErr(c_error.ErrIllegalAccess)
	}
	policyKey := fmt.Sprintf(consts.LoginPolicyCacheKeyFmt, user.LoginName, policyParam.SessionId)

	redisClient := cache.GetRedisClient()
	// 获取登录策略
	val, err := redisClient.Get(_self.Context, policyKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return base.ResultFailureErr(c_error.ErrIllegalAccess)
		}
		return base.ResultFailureErr(err)
	}
	if len(val) < 0 {
		return base.ResultFailureErr(c_error.ErrIllegalAccess)
	}
	var loginPolicy dto.LoginPolicyResultDTO
	if err = json.Unmarshal([]byte(val), &loginPolicy); err != nil {
		return base.ResultFailureErr(err)
	}
	operate := policyParam.Operate
	deleteIdx := -1
	for idx, code := range loginPolicy.PolicyCode {
		if code == policyParam.PolicyCode {
			deleteIdx = idx
			break
		}
	}

	if deleteIdx == -1 {
		return base.ResultFailureErr(c_error.ErrIllegalAccess)
	}

	var sysLoginPolicyService SysLoginPolicy
	_self.MakeService(&sysLoginPolicyService)
	//处理策略
	user.SessionId = loginDTO.PolicyParam.SessionId
	res = sysLoginPolicyService.ProcessLoginPolicy(loginDTO.PolicyParam, user, policyDTOMap)

	if !res.IsSuccess() {
		return
	}

	if operate != 1 {
		res = base.ResultCustom(203, res.Data, res.Message)
		return

	}
	// 校验成功 删除已校验策略
	loginPolicy.PolicyCode = append(loginPolicy.PolicyCode[:deleteIdx], loginPolicy.PolicyCode[deleteIdx+1:]...)
	//  Mode-->0:通过全部 or 通过一个即可
	// 没有剩余的策略
	// 通过
	if loginPolicy.Mode == "1" || len(loginPolicy.PolicyCode) <= 0 {
		// 登录成功 返回nil
		return nil
	}
	// 更新缓存
	expirationTime, err := redisClient.TTL(_self.Context, policyKey).Result()
	if err != nil {
		return base.ResultFailureErr(err)
	}
	// 更新缓存loginPolicy
	marshal, err := json.Marshal(loginPolicy)
	if err != nil {
		return base.ResultFailureErr(err)
	}
	if err = redisClient.Set(_self.Context, policyKey, marshal, expirationTime).Err(); err != nil {
		return base.ResultFailureErr(err)
	}
	// 继续返回策略信息
	res = base.ResultCustom(201, loginPolicy, res.Message+",请继续填写登录信息")

	return
}
