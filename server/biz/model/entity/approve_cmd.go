package entity

import "go-protector/server/internal/consts/table_name"

type ApproveCmd struct {
	ModelId
	AssetId uint64 `gorm:"index;comment:资产ID"  json:"assetId" binding:"required"`
	Cmd     string `gorm:"size:64;comment:审批指令"  json:"cmd" binding:"required"`
	ModelControl
	ModelDelete
}

func (_self ApproveCmd) TableName() string {
	return table_name.ApproveCmd

}
