package entity

import "go-protector/server/internal/consts/table_name"

type SysDictType struct {
	ModelId
	TypeCode string `json:"typeCode" gorm:"size:32;comment:类型编码"`
	TypeName string `json:"typeName" gorm:"size:32;comment:类型名称"`
	ModelControl
	ModelDelete
}

func (_self SysDictType) TableName() string {
	return table_name.SysDictType
}
