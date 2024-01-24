package c_logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/core/consts"
	"go-protector/server/core/utils"
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

func GetLogger(c *gin.Context) (log *SelfLogger) {
	var ok bool
	var value any
	if value, ok = c.Get(consts.CtxKeyLog); ok {
		if log, ok = value.(*SelfLogger); ok {
			return
		}
		return
	}
	var traceId string
	if value, ok = c.Get(consts.CtxKeyTraceId); !ok {
		traceId = strconv.FormatUint(utils.GetNextId(), 10)
	} else {
		if traceId, ok = value.(string); !ok {
			traceId = strconv.FormatUint(utils.GetNextId(), 10)
		}
	}
	c.Set(consts.CtxKeyTraceId, traceId)
	log = &SelfLogger{
		zapLog: _log.zapLog.Named("traceId: " + traceId),
		ctx:    c,
	}
	c.Set(consts.CtxKeyLog, log)
	return
}
func GetLoggerByCtx(ctx context.Context) (log *SelfLogger) {
	var ok bool
	log, ok = ctx.Value(consts.CtxKeyLog).(*SelfLogger)
	if ok {
		return
	}
	log = &SelfLogger{
		zapLog: _log.zapLog.With(),
		ctx:    ctx,
	}
	if traceId, ok := ctx.Value(consts.CtxKeyTraceId).(string); ok {
		log.zapLog.Named(traceId)
	}

	log.ctx = context.WithValue(ctx, consts.CtxKeyLog, log)
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
