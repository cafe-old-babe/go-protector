package dto

import (
	"go-protector/server/internal/current"
	"time"
)

type LoginDTO struct {
	// 3-7	【实战】Gin如何优雅解决数据的绑定校验
	LoginName string `json:"loginName"  binding:"required"`
	Password  string `json:"password"  binding:"required_without=PolicyParam"`
	// region 验证码相关 // 3-6	【实战】登录实现静态密码+图片验证码（使用Redis存储验证码信息）
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
