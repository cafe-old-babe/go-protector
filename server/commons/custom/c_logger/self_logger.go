package c_logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/local"
	"go-protector/server/commons/utils"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

var _log SelfLogger

type SelfLogger struct {
	zapLog *zap.Logger
	ctx    context.Context
}

var once sync.Once

func SetLogger(logger *zap.Logger) {
	once.Do(func() {
		_log.zapLog = logger
		zap.ReplaceGlobals(_log.zapLog)
	})

}

func GetLogger(ctx *gin.Context) (log *SelfLogger) {
	var ok bool
	var value any
	if value, ok = ctx.Get(local.CtxKeyLog); ok {
		if log, ok = value.(*SelfLogger); ok {
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
	log = &SelfLogger{
		zapLog: _log.zapLog.Named("traceId: " + traceId),
		ctx:    ctx,
	}

	//log = NewLoggerByOpt(ctx, zap.Fields(zapcore.Field{
	//	Key:    "traceId",
	//	Type:   zapcore.StringType,
	//	String: traceId,
	//}))
	ctx.Set(local.CtxKeyLog, log)
	return
}

func NewLoggerByField(ctx context.Context, fields ...zap.Field) *SelfLogger {
	return &SelfLogger{
		zapLog: _log.zapLog.With(fields...),
		ctx:    ctx,
	}
}

func (_self *SelfLogger) Debug(msg string, data ...any) {
	_self.zapLog.Debug(fmt.Sprintf(msg, data...))
}

func (_self *SelfLogger) DebugZap(msg string, fields ...zap.Field) {
	_self.zapLog.Debug(msg, fields...)
}

func (_self *SelfLogger) Info(msg string, data ...any) {
	_self.zapLog.Info(fmt.Sprintf(msg, data...))
}

func (_self *SelfLogger) InfoZap(msg string, fields ...zap.Field) {
	_self.zapLog.Info(msg, fields...)
}

func (_self *SelfLogger) Warn(msg string, data ...any) {
	_self.zapLog.Warn(fmt.Sprintf(msg, data...))
}

func (_self *SelfLogger) Error(msg string, data ...any) {
	_self.zapLog.Error(fmt.Sprintf(msg, data...))
}

func (_self *SelfLogger) Fatal(msg string, data ...any) {
	_self.zapLog.Fatal(fmt.Sprintf(msg, data...))
}
func (_self *SelfLogger) Panic(msg string, data ...any) {
	_self.zapLog.Panic(fmt.Sprintf(msg, data...))
}
