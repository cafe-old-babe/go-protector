package entity

type SysOtpBind struct {
	UserId     uint64 `gorm:"comment:用户ID"  json:"-"`
	Secret     string `gorm:"comment:秘钥,size:128"  json:"-"`
	OtpAuthURL string `gorm:"size:4096;comment:二维码URL"  json:"-"`
}
