package entity

import (
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_type"
)

// AssetAuth 授权表
type AssetAuth struct {
	ModelId
	AssetId    uint64      `gorm:"comment:资产id"  json:"assetId" binding:"required"`
	AssetName  string      `gorm:"size:32;comment:资产名称"  json:"assetName" binding:"required"`
	AssetIp    string      `gorm:"size:32;comment:资产IP"  json:"assetIp" binding:"required"`
	AssetAccId uint64      `gorm:"comment:资产从帐号ID"  json:"assetAccId" binding:"required"`
	AssetAcc   string      `gorm:"size:32;comment:资产从帐号"  json:"assetAcc" binding:"required"`
	UserId     uint64      `gorm:"comment:主帐号ID"  json:"userId" binding:"required"`
	UserAcc    string      `gorm:"size:32;comment:主帐号"  json:"userAcc" binding:"required"`
	StartDate  c_type.Time `gorm:"type:date;comment:授权开始时间"  json:"startDate"  binding:"required_with=EndDate,ltfield=EndDate"`
	EndDate    c_type.Time `gorm:"type:date;comment:授权结束时间"  json:"endDate" binding:"required_with=StartDate,gtfield=StartDate"`
	ModelControl
	ModelDelete
}

func (*AssetAuth) TableName() string {
	return table_name.AssetAccount
}
