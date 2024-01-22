package service

import (
	"database/sql"
	"errors"
	"go-protector/server/core/base"
	"go-protector/server/core/consts"
	"go-protector/server/core/custom/c_captcha"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_jwt"
	"go-protector/server/core/utils"
	"go-protector/server/dao"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
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
	if err = dao.SysUser.UpdateLastLoginIp(_self.DB, entity.ID, _self.Ctx.ClientIP()); err != nil {
		_self.Logger.Error("用户: %s UpdateLastLoginIp err: %v", entity.LoginName, err)
	}

	// 生成Token
	userDTO := &dto.CurrentUser{
		ID:        entity.ID,
		SessionId: utils.GetNextId(),
		LoginName: entity.LoginName,
		UserName:  entity.Username,
		LoginTime: time.Now().Format(time.DateTime),
		LoginIp:   _self.Ctx.ClientIP(),
		Avatar:    "https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png",
	}
	jwtString, expireAt, err := c_jwt.GenerateToken(userDTO)
	tempMap := map[string]any{
		"id": "admin",
	}

	res = dto.ResultSuccess(dto.LoginSuccess{
		SysUser:     userDTO,
		Token:       jwtString,
		ExpireAt:    expireAt,
		Permissions: tempMap,
		Roles:       tempMap,
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
			LockReason: sql.NullString{Valid: true, String: dto.LockReason},
			UserStatus: dto.UserStatus,
		}); err != nil {
			_self.Logger.Error("SetStatus LockUser err: %v", err)
		}
	}

	return
}
