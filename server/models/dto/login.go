package dto

type Login struct {
	LoginName string `json:"loginName,omitempty"`
	Password  string `json:"password,omitempty"`
}

type LoginSuccess struct {
	SysUser *CurrentUser `json:"sysUser"`
	Token   string       `json:"token"`
}
