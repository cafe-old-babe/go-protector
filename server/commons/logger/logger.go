package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sync"
)

var _logger *zap.Logger

var once sync.Once

func SetLogger(logger *zap.Logger) {
	once.Do(func() {
		_logger = logger
		zap.ReplaceGlobals(_logger)
	})

}
func GetLogger() *zap.Logger {
	return _logger.WithOptions()
}

func DebugF(ctx context.Context, format string, data ...any) {
	namedLogger := _logger
	if traceId, ok := ctx.Value("traceId").(string); ok {
		namedLogger = _logger.Named(traceId)
	}
	namedLogger.Debug(fmt.Sprintf(format, data...))
}

func Debug(ctx context.Context, msg string) {
	namedLogger := _logger
	if traceId, ok := ctx.Value("traceId").(string); ok {
		namedLogger = _logger.Named(traceId)
	}
	namedLogger.Debug(msg)
}

func Info(ctx context.Context, msg string) {
	namedLogger := _logger
	if traceId, ok := ctx.Value("traceId").(string); ok {
		namedLogger = _logger.Named(traceId)
	}
	namedLogger.Info(msg)
}

func InfoF(ctx context.Context, format string, data ...any) {
	namedLogger := _logger
	if traceId, ok := ctx.Value("traceId").(string); ok {
		namedLogger = _logger.Named(traceId)
	}
	namedLogger.Info(fmt.Sprintf(format, data...))
}
