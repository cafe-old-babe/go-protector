package dto

import (
	"database/sql"
)

type FindUserDTO struct {
	Id          uint64
	LoginName   string
	UserStatus  int // 如果非零 查询所有状态
	IsUnscoped  bool
	CurrentTime sql.NullTime
}

type SysUser struct {
	LoginName     string `json:"loginName"`
	UserName      string `json:"userName"`
	LastLoginTime string `json:"lastLoginTime"`
	LastLoginIp   string `json:"lastLoginIp"`
}
