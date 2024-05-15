package c_logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	log = &SelfLogger{}
	var traceId string
	if traceId, ok = ctx.Value(consts.CtxKeyTraceId).(string); ok {
		log.zapLog = _log.zapLog.Named(traceId)
	} else {
		log.zapLog = _log.zapLog.Named("temp-name")

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
	_self.zapLog.Log(zapcore.DebugLevel, msg, fields...)
}

func (_self *SelfLogger) Info(msg string, data ...any) {
	_self.zapLog.Info(fmt.Sprintf(msg, data...))
}

func (_self *SelfLogger) InfoZap(msg string, fields ...zap.Field) {
	_self.zapLog.Log(zapcore.InfoLevel, msg, fields...)
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

func Debug(msg string, data ...any) {
	_log.Debug(fmt.Sprintf(msg, data...))
}

func DebugZap(msg string, fields ...zap.Field) {
	_log.DebugZap(msg, fields...)
}

func Info(msg string, data ...any) {
	_log.Info(fmt.Sprintf(msg, data...))
}

func InfoZap(msg string, fields ...zap.Field) {
	_log.InfoZap(msg, fields...)
}

func Warn(msg string, data ...any) {
	_log.Warn(fmt.Sprintf(msg, data...))
}

func Error(msg string, data ...any) {
	_log.Error(fmt.Sprintf(msg, data...))
}
func ErrorZap(msg string, fields ...zap.Field) {
	_log.zapLog.Error(msg, fields...)
}

func Fatal(msg string, data ...any) {
	_log.Fatal(fmt.Sprintf(msg, data...))
}
func Panic(msg string, data ...any) {
	_log.Panic(fmt.Sprintf(msg, data...))
}
