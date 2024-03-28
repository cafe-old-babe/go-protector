package c_error

import "errors"

var (
	// ErrParamInvalid 参数校验不通过
	ErrParamInvalid                 = errors.New("参数校验未通过")
	ErrRecordNotFoundSysUser        = errors.New("用户不存或已锁定")
	ErrNotFoundRoleOfSysUser        = errors.New("用户未设置角色,请联系管理员")
	ErrLoginNameOrPasswordIncorrect = errors.New("用户名或密码不正确")
	ErrUpdateFailure                = errors.New("更新失败")
	ErrDeleteFailure                = errors.New("删除失败")
	ErrAuthFailure                  = errors.New("认证失败，请重新登录")
	ErrIllegalAccess                = errors.New("非法访问")
)
