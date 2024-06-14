package entity

import (
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/utils/gm"
	"gorm.io/gorm"
)

type SsoSession struct {
	ModelId
	AuthId         uint64               `gorm:"comment:授权ID"  json:"authId" binding:"required"`
	AssetId        uint64               `gorm:"comment:资产id"  json:"assetId" binding:"required"`
	AssetName      string               `gorm:"size:32;comment:资产名称" json:"assetName" binding:"required"  `
	AssetIp        string               `gorm:"size:32;comment:资产IP" json:"assetIp" binding:"required"`
	AssetPort      int                  `gorm:"size:32;comment:资产端口" json:"assetPort" binding:"required"`
	AssetGatewayId uint64               `gorm:"size:32;comment:资产网关" json:"assetGatewayId"`
	AssetAccId     uint64               `gorm:"comment:资产从帐号ID"  json:"assetAccId" binding:"required"`
	AssetAcc       string               `gorm:"size:32;comment:资产从帐号" json:"assetAcc" binding:"required"`
	AssetAccPwd    string               `gorm:"size:256;comment:资产从帐号密码" json:"assetAccPwd" binding:"required"`
	UserId         uint64               `gorm:"comment:主帐号ID"  json:"userId" binding:"required"`
	UserAcc        string               `gorm:"size:32;comment:主帐号"  json:"userAcc" binding:"required"`
	ConnectAt      c_type.Time          `gorm:"comment:连接时间" json:"connectAt" `
	Status         c_type.SessionStatus `gorm:"size:1;comment:会话状态,0-等待连接,1-正在连接,2-已连接,3-会话结束" json:"status"  binding:"required"`
	StatusText     string               `gorm:"-" json:"statusText"`
	ModelControl
	ModelDelete
}

func (_self *SsoSession) TableName() string {
	return table_name.SsoSession
}

func (_self *SsoSession) AfterFind(db *gorm.DB) (err error) {
	if len(_self.AssetAccPwd) > 0 {
		_self.AssetAccPwd, err = gm.Sm4DecryptCBC(_self.AssetAccPwd)
	}
	_self.Completion()
	return
}

func (_self *SsoSession) Completion() {
	switch _self.Status {
	case "0":
		_self.StatusText = "等待连接"
	case "1":
		_self.StatusText = "连接中"
	case "2":
		_self.StatusText = "已连接"
	case "3":
		_self.StatusText = "会话结束"

	}
}

func (_self *SsoSession) BeforeSave(db *gorm.DB) (err error) {
	if len(_self.AssetAccPwd) > 0 {
		// 先解密,如果解密失败,代表是明文
		if _, err = gm.Sm4DecryptCBC(_self.AssetAccPwd); err != nil {
			_self.AssetAccPwd, err = gm.Sm4EncryptCBC(_self.AssetAccPwd)
		}
	}
	return
}
