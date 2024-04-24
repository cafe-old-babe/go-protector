package entity

import "go-protector/server/internal/consts/table_name"

type AssetGateway struct {
	ModelId
	AgName string `gorm:"size:64;comment:网关名称" json:"agName,omitempty" binding:"required"`
	//agType     string `gorm:"size:64;comment:网络域类型" json:"agType,omitempty" binding:"required"`
	AgIp       string `gorm:"size:64;comment:网关IP" json:"agIp,omitempty" binding:"required,ip"`
	AgPort     int    `gorm:"size:64;comment:网关端口" json:"agPort,omitempty" binding:"required,min=1,max=65535"`
	AgAccount  string `gorm:"size:64;comment:网关使用帐号" json:"agAccount,omitempty" `
	AgPassword string `gorm:"size:64;comment:网关使用密码" json:"agPassword,omitempty"`
	ModelControl
	ModelDelete
}

func (AssetGateway) TableName() string {
	return table_name.AssetGateway
}
