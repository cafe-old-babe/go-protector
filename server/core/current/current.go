package current

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-protector/server/core/consts"
	"go-protector/server/models/dto"
)

func GetUserId(c context.Context) uint64 {
	if data, ok := c.Value(consts.CtxKeyUserId).(uint64); ok {
		return data
	}
	return 0
}

func SetUserId(ctx *gin.Context, data uint64) {
	if data <= 0 {
		return
	}
	ctx.Set(consts.CtxKeyUserId, data)
}

// SetUser 设置当前用户
func SetUser(c context.Context, user *dto.CurrentUser) (nc context.Context) {
	if c == nil || user == nil {
		return
	}
	nc = context.WithValue(c, consts.CtxKeyCurrentUser, user)

	return
}

// GetUser 获取当前用户
func GetUser(c context.Context) (user *dto.CurrentUser, ok bool) {
	user, ok = c.Value(c).(*dto.CurrentUser)
	return
}
