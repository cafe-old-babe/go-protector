package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-protector/server/core/base"
	"go-protector/server/core/consts"
	"go-protector/server/core/consts/table_name"
	"go-protector/server/core/current"
	"go-protector/server/core/custom/c_captcha"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_jwt"
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

// DoLogin 登录验证
func (_self *SysUser) DoLogin(loginDTO dto.Login) (res *dto.Result) {

	if !c_captcha.Verify(loginDTO.Cid, loginDTO.Code, true) {
		return dto.ResultFailureMsg("验证码错误或已失效")
	}

	var sysUser *entity.SysUser
	var err error
	sysUser, err = dao.SysUser.FindUserByDTO(_self.DB, &dto.FindUser{
		LoginName: loginDTO.LoginName,
	})

	if err != nil {
		_self.Logger.Error("FindUser: %s, err: %v", loginDTO.LoginName, err.Error())
		if errors.Is(c_error.ErrRecordNotFoundSysUser, err) {
			res = dto.ResultFailureMsg(c_error.ErrLoginNameOrPasswordIncorrect.Error())
		}
		res = dto.ResultFailureMsg(err.Error())
		return
	}

	if sysUser == nil {
		_self.Logger.Error("未查询到用户: %s", loginDTO.LoginName)
		res = dto.ResultFailureMsg(c_error.ErrLoginNameOrPasswordIncorrect.Error())
		return
	}
	now := time.Now()
	// 检查有效期
	if sysUser.ExpirationAt.Valid {

		if now.After(sysUser.ExpirationAt.Time) {
			_self.Logger.Error("用户: %s 已过有效期", loginDTO.LoginName)
			res = dto.ResultFailureMsg(c_error.ErrLoginNameOrPasswordIncorrect.Error())
			// 更新用户信息
			sysUser.UserStatus = consts.LockTypeExpire
			sysUser.LockReason = sql.NullString{
				String: "用户已过有效期",
				Valid:  true,
			}
			sysUser.UpdatedBy = sysUser.ID
			if err = dao.SysUser.LockUser(_self.DB, sysUser); err != nil {
				_self.Logger.Error("用户: %s lockUser err: %v", loginDTO.LoginName, err)
			}
			return
		}
	}
	// 校验密码
	if sysUser.Password != loginDTO.Password {
		res = dto.ResultFailureMsg(c_error.ErrLoginNameOrPasswordIncorrect.Error())
		return
	}

	return _self.LoginSuccess(sysUser)

}

// LoginSuccess 登录成功
func (_self *SysUser) LoginSuccess(entity *entity.SysUser) (res *dto.Result) {
	var err error
	// 更新最后登录时间 最后登录IP
	if err = dao.SysUser.UpdateLastLoginIp(_self.DB, entity.ID, _self.Context.ClientIP()); err != nil {
		_self.Logger.Error("用户: %s UpdateLastLoginIp err: %v", entity.LoginName, err)
	}

	// 生成Token
	userDTO := &dto.CurrentUser{
		ID:        entity.ID,
		SessionId: utils.GetNextId(),
		LoginName: entity.LoginName,
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

	jwtString, expireAt, err := c_jwt.GenerateToken(userDTO)

	res = dto.ResultSuccess(dto.LoginSuccess{
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
func (_self *SysUser) UserInfo() (res *dto.Result) {
	// 查询用户角色
	user, ok := current.GetUser(_self.Context)
	if !ok {
		res = dto.ResultFailureMsg("获取当前用户失败")
		return
	}
	// 查询角色关联菜单
	roleIds := user.RoleIds
	if len(roleIds) <= 0 {
		res = dto.ResultFailureMsg("当前用户未绑定角色")
		return
	}
	var roleService SysRole
	_self.MakeService(&roleService)
	menuSlice, buttonSlice, err := roleService.GetMenuByRoleIds(roleIds, user.IsAdmin)
	if err != nil {
		return dto.ResultFailureErr(err)
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
	return dto.ResultSuccess(roleInfo)

}

// Nav 获取菜单
func (_self *SysUser) Nav() (res *dto.Result) {
	// 查询用户角色
	user, ok := current.GetUser(_self.Context)
	if !ok {
		res = dto.ResultFailureMsg("获取当前用户失败")
		return
	}
	// 查询角色关联菜单
	roleIds := user.RoleIds
	if len(roleIds) <= 0 {
		res = dto.ResultFailureMsg("当前用户未绑定角色")
		return
	}

	menuSlice, err := dao.SysRole.GetPermissionSliceByIds(_self.DB, roleIds, []int8{0, 1}, user.IsAdmin)
	if err != nil {
		res = dto.ResultFailureErr(err)
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
	return dto.ResultSuccess(menuInfoSlice)
}

// Page 人员分页查询
func (_self *SysUser) Page(req *dto.UserPageReq) (result *dto.Result) {
	if req == nil {
		result = dto.ResultFailureErr(c_error.ErrParamInvalid)
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
		result = dto.ResultFailureErr(err)
	} else {
		result = dto.ResultPage(voSlice, req.GetPagination(), count)
	}
	return
}

func (_self *SysUser) Save(req *dto.UserSaveReq) (result *dto.Result) {

	if err := dao.SysUser.Save(_self.DB, req); err != nil {
		return dto.ResultFailureErr(err)
	}

	return dto.ResultSuccessMsg("保存成功")

}

func (_self *SysUser) DeleteByIds(req *dto.IdsReq) (result *dto.Result) {

	if req == nil || len(req.Ids) <= 0 {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
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
		result = dto.ResultFailureErr(err)
		return
	}

	// delete post_relation
	if err := db.Delete(&entity.SysPostRelation{},
		"relation_id in (?) and relation_type = ?", req.Ids, consts.User).Error; err != nil {
		result = dto.ResultFailureErr(err)
		return
	}

	// delete role_relation
	if err := db.Delete(&entity.SysRoleRelation{},
		"relation_id in (?) and relation_type = ?", req.Ids, consts.User).Error; err != nil {
		result = dto.ResultFailureErr(err)
		return
	}

	return dto.ResultSuccessMsg("删除成功")
}
