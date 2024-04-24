package entity

import "go-protector/server/internal/consts/table_name"

type AssetGroup struct {
	ModelId
	Name string `gorm:"size:32;comment:资源组名称"  json:"name" binding:"required" `
	PID  uint64 `gorm:"comment:父级ID" json:"pid" binding:"min=0"`
	Sort int32  `gorm:"comment:排序字段" json:"sort"`
	ModelControl
	ModelDelete
}

func (_self AssetGroup) TableName() string {
	return table_name.AssetGroup
}
