package service

import (
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_captcha"
	"go-protector/server/internal/custom/c_error"
)

// DoLogin 登录验证
func (_self *SysUser) DoLogin(loginDTO dto.LoginDTO) (res *base.Result) {

	if loginDTO.PolicyParam == nil {
		if !c_captcha.Verify(loginDTO.Cid, loginDTO.Code, true) {
			return base.ResultFailureMsg("验证码错误或已失效")
		}
	}

	res = _self.FindUserByDTO(&dto.FindUserDTO{
		LoginName: loginDTO.LoginName,
	})
	if !res.IsSuccess() {
		_self.GetLogger().Error("FindUser: %s, err: %v", loginDTO.LoginName, res.Message)
		return base.ResultFailureErr(c_error.ErrLoginNameOrPasswordIncorrect)
	}
	sysUser := res.Data.(*entity.SysUser)

	return _self.checkLogin(&loginDTO, sysUser)

	//return _self.LoginSuccess(sysUser)

}
