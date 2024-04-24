package entity

import (
	"go-protector/server/internal/consts/table_name"
)

type AssetBasic struct {
	ModelId
	AssetName      string `gorm:"size:32;comment:资产名称"  json:"assetName" binding:"required" `
	AssetGroupId   uint64 `gorm:"comment:资源组ID"  json:"groupId" binding:"required"`
	IP             string `gorm:"size:32;comment:ip"  json:"ip" binding:"required,ip"`
	Port           uint   `gorm:"comment:端口"  json:"port" binding:"required,min=1,max=65535" `
	ManagerUserId  uint64 `gorm:"comment:资产管理员Id"  json:"managerUserId" binding:"required"`
	AssetGatewayId uint64 `gorm:"comment:网关ID"  json:"agId"`
	ModelControl
	ModelDelete
}

func (_self AssetBasic) TableName() string {
	return table_name.AssetBasic
}

// AssetInfo belongs to
type AssetInfo struct {
	AssetBasic //`gorm:"embedded"`
	//RootAcc      string       `json:"rootAcc"`
	//RootPwd      string       `json:"-"`
	RootAcc      AssetAccount `gorm:"foreignKey:ID;references:AssetId" json:"rootAcc"`
	AssetGroup   AssetGroup   `json:"assetGroup"`
	ManagerUser  SysUser      `json:"managerUser"`
	AssetGateway AssetGateway `gorm:"foreignKey:AssetGatewayId" json:"assetGateway"` // belongs to
}

type AssetInfoAccount struct {
	AssetInfo
	Accounts []AssetAccount `gorm:"foreignKey:AssetId;references:ID" json:"accounts"`
}
