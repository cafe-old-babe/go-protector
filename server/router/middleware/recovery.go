package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/logger"
	"go-protector/server/commons/result"
)

// Recovery 全局 recover
func Recovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.NewLogger(ctx).Error("recover err: %v", err)
			if ctx.IsAborted() {
				ctx.Status(200)
			}
			result.Failure(ctx, nil, fmt.Sprintf("recover err: %v", err))
		}
		return
	}()
	ctx.Next()

}
