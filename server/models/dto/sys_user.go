package dto

import (
	"database/sql"
)

type FindUser struct {
	ID          uint64
	LoginName   string
	UserStatus  int // 如果非零 查询所有状态
	IsUnscoped  bool
	CurrentTime sql.NullTime
}

type CurrentUser struct {
	ID        uint64 `json:"ID"`
	SessionId uint64 `json:"sessionId"`
	LoginName string `json:"loginName"`
	UserName  string `json:"userName"`
	LoginTime string `json:"loginTime"`
	LoginIp   string `json:"loginIp"`
}

type SetStatus struct {
	ID           uint64
	UserStatus   int
	LockReason   string
	ExpirationAt string
}
