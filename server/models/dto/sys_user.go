package dto

import (
	"database/sql"
	"go-protector/server/core/custom/c_type"
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
	IsAdmin   bool     `json:"isAdmin"`
}

type SetStatus struct {
	ID           uint64
	UserStatus   int
	LockReason   string
	ExpirationAt c_type.Time
}

// UserPageReq 人员管理分页查询
type UserPageReq struct {
	Pagination
	DeptIds   []uint64 `json:"deptIds"`
	LoginName string   `json:"loginName"`
	Username  string   `json:"username"`
}

type UserSaveReq struct {
	DeptId       uint64      `json:"deptId" binding:"required"`
	PostIds      []uint64    `json:"postIds"  binding:"required"`
	RoleIds      []uint64    `json:"roleIds"  binding:"required"`
	LoginName    string      `json:"loginName"  binding:"required"`
	Email        string      `json:"email"  binding:"required,email"`
	Password     string      `json:"password"  binding:"required_without=ID"`
	Username     string      `json:"username"  binding:"required"`
	Sex          string      `json:"sex"  binding:"required"`
	ExpirationAt c_type.Time `json:"expirationAt"`
	ID           uint64      `json:"id"`
}
