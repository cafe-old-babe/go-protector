package dto

import "time"

type Login struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
	// region 验证码相关
	Cid  string `json:"cid"`
	Code string `json:"code"`
	// endregion
}

type LoginSuccess struct {
	SysUser     *CurrentUser   `json:"user"`
	Token       string         `json:"token"`
	ExpireAt    time.Time      `json:"expireAt"`
	Permissions map[string]any `json:"permissions"`
	Roles       map[string]any `json:"roles"`
}
