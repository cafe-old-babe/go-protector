package dto

import (
	"go-protector/server/core/current"
	"time"
)

type LoginDTO struct {
	LoginName string `json:"loginName"  binding:"required"`
	Password  string `json:"password"  binding:"required"`
	// region 验证码相关
	Cid  string `json:"cid"`
	Code string `json:"code"`
	// endregion
	PolicyParam *LoginPolicyParamDTO `json:"policyParam"`
}

type LoginSuccess struct {
	SysUser     *current.User  `json:"user"`
	Token       string         `json:"token"`
	ExpireAt    time.Time      `json:"expireAt"`
	Permissions map[string]any `json:"permissions"`
	Roles       map[string]any `json:"roles"`
}
