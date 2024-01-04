package service

import (
	"database/sql"
	"errors"
	"go-protector/server/commons/base"
	"go-protector/server/commons/custom/c_error"
	"go-protector/server/commons/local"
	"go-protector/server/dao"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"time"
)

type SysUser struct {
	base.Service
}

func (_self *SysUser) DoLogin(loginDTO dto.Login) (res *dto.Result) {
	var sysUser *entity.SysUser
	var err error
	sysUser, err = dao.SysUser.FindUserByDTO(_self.DB, &dto.FindUserDTO{
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
			sysUser.UserStatus = local.LockTypeExpire
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
func (_self *SysUser) LoginSuccess(sysUser *entity.SysUser) (res *dto.Result) {
	var err error
	// 更新最后登录时间 最后登录IP
	if err = dao.SysUser.UpdateLastLoginIp(_self.DB, sysUser.ID, _self.Ctx.ClientIP()); err != nil {
		_self.Logger.Error("用户: %s UpdateLastLoginIp err: %v", sysUser.LoginName, err)
	}
	res = dto.ResultSuccess(dto.LoginSuccess{
		SysUser: &dto.SysUser{
			LoginName:     sysUser.LoginName,
			UserName:      sysUser.Username,
			LastLoginTime: time.Now().Format(time.DateTime),
			LastLoginIp:   _self.Ctx.ClientIP(),
		},
		// todo jwt token
		Token: "",
	})
	return
}
