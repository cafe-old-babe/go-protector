package entity

import "go-protector/server/core/consts/table_name"

type AssetNetwork struct {
	ModelId
	AnName string `gorm:"size:64;comment:网络域名称" json:"anName,omitempty" binding:"required"`
	//AnType     string `gorm:"size:64;comment:网络域类型" json:"anType,omitempty" binding:"required"`
	AnIp       string `gorm:"size:64;comment:网络域IP" json:"anIp,omitempty" binding:"required,ip"`
	AnPort     int    `gorm:"size:64;comment:网络域端口" json:"anPort,omitempty" binding:"required,min=1,max=65535"`
	AnAccount  string `gorm:"size:64;comment:网络域使用帐号" json:"anAccount,omitempty" `
	AnPassword string `gorm:"size:64;comment:网络域使用密码" json:"anPassword,omitempty"`
	ModelControl
	ModelDelete
}

func (AssetNetwork) TableName() string {
	return table_name.AssetNetwork
}
