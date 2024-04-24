package entity

import "go-protector/server/internal/consts/table_name"

type SysOtpBind struct {
	UserId     uint64 `gorm:"comment:用户ID;primaryKey"  json:"-"`
	OtpAuthURL string `gorm:"size:4096;comment:二维码URL"  json:"-"`
}

func (SysOtpBind) TableName() string {

	return table_name.SysOtpBind
}
