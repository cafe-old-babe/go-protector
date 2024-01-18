package current

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-protector/server/core/consts"
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
