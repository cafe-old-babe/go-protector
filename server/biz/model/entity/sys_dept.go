package entity

import "go-protector/server/internal/consts/table_name"

type SysDept struct {
	ModelId
	DeptName string `gorm:"size:32;comment:部门名称"  json:"deptName" binding:"required" `
	PID      uint64 `gorm:"comment:父级ID" json:"pid" binding:"min=0"`
	PathName string `gorm:"size:32;comment:部门全路径" json:"pathName"`
	Sort     int32  `gorm:"comment:排序字段" json:"sort"`
	ModelControl
	ModelDelete
}

func (_self SysDept) TableName() string {
	return table_name.SysDept
}
