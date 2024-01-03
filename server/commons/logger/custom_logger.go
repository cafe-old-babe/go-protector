package logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/local"
	"go-protector/server/commons/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"sync"
)

var log CustomLogger

type CustomLogger struct {
	zapLog *zap.Logger
	ctx    context.Context
}

var once sync.Once

func SetLogger(logger *zap.Logger) {
	once.Do(func() {
		log.zapLog = logger
		zap.ReplaceGlobals(log.zapLog)
	})

}

func NewLogger(ctx *gin.Context) (log *CustomLogger) {
	var ok bool
	var value any
	if value, ok = ctx.Get(local.CtxKeyLog); ok {
		if _, ok = value.(*CustomLogger); ok {
			return
		}
		return
	}
	var traceId string
	if value, ok = ctx.Get(local.CtxKeyTraceId); !ok {
		traceId = strconv.FormatUint(utils.GetNextId(), 10)
	} else {
		if traceId, ok = value.(string); !ok {
			traceId = strconv.FormatUint(utils.GetNextId(), 10)
		}
	}
	ctx.Set(local.CtxKeyTraceId, traceId)

	log = NewLoggerByOpt(ctx, zap.Fields(zapcore.Field{
		Key:    "traceId",
		Type:   zapcore.StringType,
		String: traceId,
	}))
	ctx.Set(local.CtxKeyLog, log)
	return
}

func NewLoggerByOpt(ctx context.Context, opt ...zap.Option) *CustomLogger {
	return &CustomLogger{
		zapLog: log.zapLog.WithOptions(opt...),
		ctx:    ctx,
	}
}

func NewLoggerByField(ctx context.Context, fields ...zap.Field) *CustomLogger {
	return &CustomLogger{
		zapLog: log.zapLog.With(fields...),
		ctx:    ctx,
	}
}

func (_self *CustomLogger) Debug(msg string, data ...any) {
	_self.zapLog.Debug(fmt.Sprintf(msg, data...))
}

func (_self *CustomLogger) DebugZap(msg string, fields ...zap.Field) {
	_self.zapLog.Debug(msg, fields...)
}

func (_self *CustomLogger) Info(msg string, data ...any) {
	_self.zapLog.Info(fmt.Sprintf(msg, data...))
}

func (_self *CustomLogger) InfoZap(msg string, fields ...zap.Field) {
	_self.zapLog.Info(msg, fields...)
}

func (_self *CustomLogger) Warn(msg string, data ...any) {
	_self.zapLog.Warn(fmt.Sprintf(msg, data...))
}

func (_self *CustomLogger) Error(msg string, data ...any) {
	_self.zapLog.Error(fmt.Sprintf(msg, data...))
}

func (_self *CustomLogger) Fatal(msg string, data ...any) {
	_self.zapLog.Fatal(fmt.Sprintf(msg, data...))
}
func (_self *CustomLogger) Panic(msg string, data ...any) {
	_self.zapLog.Panic(fmt.Sprintf(msg, data...))
}
