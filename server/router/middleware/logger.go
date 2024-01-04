package middleware

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/custom/c_logger"
	"go.uber.org/zap"
	"time"
)

// RecordLog 记录日志
func RecordLog(ctx *gin.Context) {
	selfLogger := c_logger.GetLogger(ctx)
	start := time.Now()
	path := ctx.Request.URL.Path
	query := ctx.Request.URL.RawQuery
	ctx.Next()

	cost := time.Since(start)
	selfLogger.DebugZap(path+" => ",
		zap.Int("status", ctx.Writer.Status()),
		zap.String("method", ctx.Request.Method),
		zap.String("path", path),
		zap.String("query", query),
		zap.String("ip", ctx.ClientIP()),
		zap.String("user-agent", ctx.Request.UserAgent()),
		zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
		zap.String("cost", cost.String()),
	)
}
