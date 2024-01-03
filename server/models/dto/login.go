package dto

type Login struct {
	LoginName string `json:"loginName,omitempty"`
	Password  string `json:"password,omitempty"`
}

type LoginSuccess struct {
	SysUser *SysUser `json:"sysUser"`
	Token   string   `json:"token"`
}
