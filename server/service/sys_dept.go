package service

import (
	"errors"
	"go-protector/server/core/base"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"gorm.io/gorm"
)

type SysDept struct {
	base.Service
}

func (_self *SysDept) DeptTree() *dto.Result {
	var deptSlice []entity.SysDept
	if err := _self.DB.Find(&deptSlice).Error; err != nil {
		return dto.ResultFailureErr(err)
	}
	node := dto.GenerateTree(deptSlice, 0, "ID", "PID", "DeptName")
	return dto.ResultSuccess(node)
}

func (_self *SysDept) SaveCheck(entity *entity.SysDept) error {
	var count int64
	err := _self.DB.Model(entity).Scopes(func(db *gorm.DB) *gorm.DB {
		if entity.ID > 0 {
			db = db.Where("id <> ?", entity.ID)
		}

		return db.Where("dept_name = ? ", entity.DeptName)
	}).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("同级别部门中不能同时存在:" + entity.DeptName)
	}
	return nil
}
