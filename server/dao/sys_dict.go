package dao

import (
	"go-protector/server/models/entity"
	"gorm.io/gorm"
)

var SysDict sysDict

type sysDict struct {
}

func (_self sysDict) DeleteTypeByIds(db *gorm.DB, ids []uint64) (err error) {
	return db.Delete(&entity.SysDictType{}, ids).Error
}

func (_self sysDict) DeleteDataByTypeCode(db *gorm.DB, typeCodeSlice []string) (err error) {
	return db.Delete(&entity.SysDictData{}, "type_code in ?", typeCodeSlice).Error
}
