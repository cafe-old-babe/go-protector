package initialize

import (
	"go-protector/server/internal/config"
	"go-protector/server/internal/consts"
	gormLogger "go-protector/server/internal/custom/c_logger/gorm"
	"go-protector/server/internal/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func initDB() error {
	dbConfig := config.GetConfig().Database
	serverCfg := config.GetConfig().Server
	dsn, err := dbConfig.GetDsn()
	if err != nil {
		return err
	}
	//var gormLog logger.Interface
	//if serverCfg.Env == "dev" {
	//	gormLog = logger.Default.LogMode(logger.Info)
	//} else {
	//	gormLog = logger.Default.LogMode(logger.Error)
	//}

	var logLevel logger.LogLevel
	var colorful bool
	if serverCfg.Env == consts.CfgEnvDebug {
		logLevel = logger.Info
		colorful = true
	} else {
		logLevel = logger.Error

	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //禁用 自动创建数据库外键约束
		Logger: gormLogger.NewGormLogger(logger.Config{
			SlowThreshold: time.Second,
			Colorful:      colorful,
			LogLevel:      logLevel,
		}),
	})
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	if dbConfig.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	}
	if dbConfig.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	}
	if dbConfig.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime) * time.Second)
	}
	if dbConfig.ConnMaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifeTime) * time.Second)
	}

	database.SetDB(db)

	return nil
}
