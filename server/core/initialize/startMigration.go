package initialize

import (
	"context"
	"go-protector/server/core/database"
	"go-protector/server/models/entity"
)

// StartMigration 启动自动迁移 创建数据库 go-protector
func StartMigration() (err error) {
	//if err = initLogger(); err != nil {
	//	return
	//}
	//if err = initDB(); err != nil {
	//	return err
	//}
	// https://gorm.io/zh_CN/docs/migration.html#AutoMigrate
	err = database.GetDB(context.Background()).
		Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(
			&entity.SysUser{},
			&entity.SysDictData{},
			&entity.SysDictType{},
			&entity.SysDept{},
			&entity.SysMenu{},
			&entity.SysRole{},
			&entity.SysRoleRelation{},
			&entity.SysPost{},
			&entity.SysPostRelation{},
			&entity.SysLoginPolicy{},
		)
	return err
}
