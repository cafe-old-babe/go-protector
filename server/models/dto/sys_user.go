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
	ID        uint64   `json:"id"`
	SessionId uint64   `json:"sessionId"`
	LoginName string   `json:"loginName"`
	UserName  string   `json:"userName"`
	LoginTime string   `json:"loginTime"`
	LoginIp   string   `json:"loginIp"`
	Avatar    string   `json:"avatar"`
	RoleIds   []uint64 `json:"roleIds"`
	DeptId    uint64   `json:"deptId"`
	IsAdmin   bool     `json:"is"`
}

type SetStatus struct {
	ID           uint64
	UserStatus   int
	LockReason   string
	ExpirationAt string
}

// UserPageReq 人员管理分页查询
type UserPageReq struct {
	Pagination
	DeptIds   []uint64 `json:"deptIds"`
	LoginName string   `json:"loginName"`
	Username  string   `json:"username"`
}
