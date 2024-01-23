package entity

import "go-protector/server/core/consts/table_name"

type SysDictData struct {
	ModelId
	TypeCode   string `gorm:"size:32;comment:类型编码" json:"typeCode"`
	DataCode   string `gorm:"size:32;comment:数据编码" json:"dataCode"`
	DataName   string `gorm:"size:32;comment:数据名称" json:"dataName"`
	Sort       int    `gorm:"size:2;comment:排序字段" json:"sort"`
	DataStatus int    `gorm:"size:1;comment:数据状态,0:正常,1:停用" json:"dataStatus"`
	ModelControl
	ModelDelete
}

func (_self SysDictType) SysDictData() string {
	return table_name.SysDictData
}
