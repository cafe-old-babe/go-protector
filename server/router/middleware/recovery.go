package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/core/custom/c_logger"
	"go-protector/server/core/custom/c_result"
)

// Recovery 全局 recover
func Recovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c_logger.GetLogger(ctx).Error("recover err: %v", err)
			if ctx.IsAborted() {
				ctx.Status(200)
			}
			c_result.Failure(ctx, nil, fmt.Sprintf("recover err: %v", err))
		}
		return
	}()
	ctx.Next()

}
