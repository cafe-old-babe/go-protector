package entity

import "go-protector/server/core/consts/table_name"

type SysDept struct {
	ModelId
	DeptName string `gorm:"size:32;comment:部门名称"  json:"deptName"`
	PID      uint64 `gorm:"comment:父级ID" json:"pid"`
	PathName string `gorm:"size:32;comment:部门全路径" json:"pathName"`
	Sort     int32  `gorm:"comment:排序字段" json:"sort"`
	Status   int8   `gorm:"comment:数据状态,0:正常,1:停用" json:"status"`
	ModelControl
	ModelDelete
}

func (_self SysDept) TableName() string {
	return table_name.SysDept
}
