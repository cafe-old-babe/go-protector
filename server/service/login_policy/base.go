package login_policy

import (
	"context"
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_type"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
)

type IPolicyService interface {
	GetPolicyService(c context.Context, user *entity.SysUser) IPolicyService
	// Process email,otp
	// operate 0:发送邮件;1:校验验证码
	// isPass 0:未通过,1:通过, 只有校验的时候才使用 isPass字段
	// err 错误信息
	Process(operate int, dto dto.ILoginPolicyDTO) (isPass int, err error)

	NextType() c_type.LoginPolicyCode
}
type commonLoginPolicyService struct {
	base.IService
	user *entity.SysUser
}

type EmailLoginPolicyService struct {
}
