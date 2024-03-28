package current

import (
	"context"
	"go-protector/server/core/consts"
)

type User struct {
	ID        uint64   `json:"id"`
	SessionId string   `json:"sessionId"`
	LoginName string   `json:"loginName"`
	Email     string   `json:"email"`
	UserName  string   `json:"userName"`
	LoginTime string   `json:"loginTime"`
	LoginIp   string   `json:"loginIp"`
	Avatar    string   `json:"avatar"`
	RoleIds   []uint64 `json:"roleIds"`
	DeptId    uint64   `json:"deptId"`
	IsAdmin   bool     `json:"isAdmin"`
}

func GetUserId(c context.Context) uint64 {
	if data, ok := c.Value(consts.CtxKeyUserId).(uint64); ok {
		return data
	}
	return 0
}

func SetUserId(c context.Context, data uint64) context.Context {
	if data <= 0 {
		return c
	}
	return context.WithValue(c, consts.CtxKeyUserId, data)
}

// SetUser 设置当前用户
func SetUser(c context.Context, user *User) (nc context.Context) {
	if c == nil || user == nil {
		return
	}
	nc = context.WithValue(SetUserId(c, user.ID), consts.CtxKeyCurrentUser, user)
	return
}

// GetUser 获取当前用户
func GetUser(c context.Context) (user *User, ok bool) {
	user, ok = c.Value(consts.CtxKeyCurrentUser).(*User)
	return
}
