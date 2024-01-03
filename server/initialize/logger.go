package initialize

import (
	"go-protector/server/commons/config"
	"go-protector/server/commons/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// initLogger 初始化日志组件
// https://www.cnblogs.com/Vikyanite/p/17210643.html
func initLogger() (err error) {
	loggerCfg := config.GetConfig().Logger
	var zapLog *zap.Logger
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.0000")

	// 返回完整调用路径
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// 将日志等级标识设置为大写并且有颜色
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 生成打印到日志文件中的encoder
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	// 将日志等级标识设置为大写并且有颜色
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// 生成打印到console的encoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	outputPath := path.Join(loggerCfg.Path, loggerCfg.FileName)
	if !path.IsAbs(loggerCfg.Path) {
		// 获取loggerCfg.Path的绝对路径
		if outputPath, err = filepath.Abs(outputPath); err != nil {
			return
		}
	}

	level := getLevel(loggerCfg.Level)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(&lumberjack.Logger{
			Filename:   outputPath,
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 30,
			LocalTime:  true,
			Compress:   true,
		}), level),
	)
	zapLog = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	logger.SetLogger(zapLog)

	return
}
func getLevel(levelStr string) zapcore.Level {
	// 统一转换小写
	levelStr = strings.ToLower(levelStr)

	switch levelStr {
	case "trace", "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}

}
