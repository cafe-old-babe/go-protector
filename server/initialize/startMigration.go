package initialize

import (
	"go-protector/server/commons/database"
	"go-protector/server/models/entity"
)

// StartMigration 启动自动迁移 创建数据库 go-protector
func StartMigration() (err error) {
	if err = initLogger(); err != nil {
		return
	}
	if err = initDB(); err != nil {
		return err
	}
	err = database.GetDB().AutoMigrate(&entity.SysUser{})
	return err
}
