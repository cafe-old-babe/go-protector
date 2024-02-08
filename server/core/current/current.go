package current

import (
	"context"
	"go-protector/server/core/consts"
	"go-protector/server/models/dto"
)

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
func SetUser(c context.Context, user *dto.CurrentUser) (nc context.Context) {
	if c == nil || user == nil {
		return
	}
	nc = context.WithValue(SetUserId(c, user.ID), consts.CtxKeyCurrentUser, user)
	return
}

// GetUser 获取当前用户
func GetUser(c context.Context) (user *dto.CurrentUser, ok bool) {
	user, ok = c.Value(consts.CtxKeyCurrentUser).(*dto.CurrentUser)
	return
}
