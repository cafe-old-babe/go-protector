package initialize

import (
	"context"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/database"
)

// startMigration 启动自动迁移 创建数据库 go-protector
func startMigration() (err error) {
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
			&entity.SysOtpBind{},
			&entity.AssetGateway{},
			&entity.AssetGroup{},
			&entity.AssetBasic{},
			&entity.AssetAccount{},
			&entity.AssetAccountExtend{},
			&entity.AssetAuth{},
		)
	return err
}
