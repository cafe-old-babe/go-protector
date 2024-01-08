package gorm_logger

import (
	"context"
	"errors"
	"fmt"
	"go-protector/server/core/custom/c_logger"
	"go-protector/server/core/local"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type gormLogger struct {
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewGormLogger(config logger.Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = logger.Green + "%s\n" + logger.Reset + logger.Green + "[info] " + logger.Reset
		warnStr = logger.BlueBold + "%s\n" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		errStr = logger.Magenta + "%s\n" + logger.Reset + logger.Red + "[error] " + logger.Reset
		traceStr = logger.Green + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		traceWarnStr = logger.Green + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}

	return &gormLogger{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

func (_self *gormLogger) getLoggerByCtx(ctx context.Context) *c_logger.SelfLogger {
	if log, ok := ctx.Value(local.CtxKeyLog).(*c_logger.SelfLogger); ok {
		return log
	}

	return c_logger.NewLoggerByField(ctx)
}

// LogMode log mode
func (_self *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	// *_self 表示将 _self 指针所指向的对象解构出来，然后将其赋值给 newLog。
	// newLog 就成为了 _self 的一个副本，具有相同的结构和数据。
	// 这种拷贝方式可以被视为一种浅拷贝。
	// 深拷贝是指复制对象的所有属性以及这些属性所引用的所有对象，以创建一个独立的副本。
	// 而这里的拷贝方式只复制了对象的属性，而没有复制属性所引用的对象。
	// 新的结构体副本会复制原有结构体中的字段值，
	// 但是如果这些字段是引用类型（例如切片、映射或指针），则新结构体和原结构体将引用相同的底层数据。
	newLog := *_self
	newLog.LogLevel = level
	return &newLog
}

// Info print info
func (_self *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if _self.LogLevel >= logger.Info {
		log := _self.getLoggerByCtx(ctx)
		log.Info(_self.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		//_self.Printf(_self.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (_self *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if _self.LogLevel >= logger.Warn {
		log := _self.getLoggerByCtx(ctx)
		log.Warn(_self.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		//_self.Printf(_self.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (_self *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if _self.LogLevel >= logger.Error {
		log := _self.getLoggerByCtx(ctx)
		log.Error(_self.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		//_self.Printf(_self.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (_self *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if _self.LogLevel <= logger.Silent {
		return
	}
	log := _self.getLoggerByCtx(ctx)
	elapsed := time.Since(begin)
	switch {
	case err != nil && _self.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !_self.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.Error(_self.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			//_self.Printf(_self.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Error(_self.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			//_self.Printf(_self.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > _self.SlowThreshold && _self.SlowThreshold != 0 && _self.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", _self.SlowThreshold)
		if rows == -1 {
			log.Warn(_self.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			//_self.Printf(_self.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Warn(_self.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			//_self.Printf(_self.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case _self.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Info(_self.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			//_self.Printf(_self.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Info(_self.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			//_self.Printf(_self.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
