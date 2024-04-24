package entity

import (
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/utils/gm"
	"gorm.io/gorm"
)

// AssetAccount 从帐号表
type AssetAccount struct {
	ModelId
	AssetId           uint64             `gorm:"comment:资产id,"  json:"assetId" binding:"required"`
	Account           string             `gorm:"size:32;comment:从帐号"  json:"account" binding:"required"`
	Password          string             `gorm:"size:256;comment:密码"  json:"password,omitempty" binding:"required_without=ID"`
	AccountType       string             `gorm:"size:1;comment:从帐号类型,-1-收集后为纳管从帐号,0-特权帐号,1-管理帐号(管理帐号可执行sudo),2-普通帐号(普通帐号不可执行sudo)"  json:"accountType" binding:"oneof=-1 0 1 2" `
	AccountTypeText   string             `gorm:"-"  json:"accountTypeText" `
	AccountStatus     string             `gorm:"size:1;comment:从帐号状态,-1-采集失败,0-未采集信息,1-正常,2-即将过期,3-已过期,4-已禁用"  json:"accountStatus" binding:"oneof=-1 0 1 2 3 4"`
	AccountStatusText string             `gorm:"-"  json:"accountStatusText"`
	Extend            AssetAccountExtend `gorm:"foreignKey:ID"  json:"extend" binding:"omitempty"`
	AssetBasic        AssetBasic         `gorm:"foreignKey:AssetId"  json:"assetBasic" binding:"omitempty"`
	ModelControl
	ModelDelete
}

func (*AssetAccount) TableName() string {
	return table_name.AssetAccount
}

func (_self *AssetAccount) BeforeSave(db *gorm.DB) (err error) {
	if len(_self.Password) > 0 {
		_self.Password, err = gm.Sm4EncryptCBC(_self.Password)
	}
	return
}

func (_self *AssetAccount) AfterCreate(db *gorm.DB) (err error) {
	if _self.AccountType == "0" {
		// 特权帐号不采集
		return
	}
	err = db.Create(AssetAccountExtend{
		ModelId: ModelId{ID: _self.ID},
	}).Error
	return
}

func (_self *AssetAccount) BeforeUpdate(db *gorm.DB) (err error) {
	if len(_self.Password) > 0 && _self.ID > 0 {
		// 更新且有密码属性 校验密文密码是否与数据库一致
		var dbPassword string
		if err = db.Model(_self).Where("id = ?", _self.ID).
			Pluck("password", &dbPassword).Error; err != nil {
			return
		}
		// 如果一样,则不加入更新
		if dbPassword == _self.Password {
			_self.Password = ""
		}
	}
	return
}

func (_self *AssetAccount) AfterFind(db *gorm.DB) (err error) {
	if len(_self.Password) > 0 {
		_self.Password, err = gm.Sm4DecryptCBC(_self.Password)
	}
	_self.completion()
	return
}

func (_self *AssetAccount) completion() {

	switch _self.AccountType {
	case "-1":
		_self.AccountTypeText = "收集从帐号"
	case "0":
		_self.AccountTypeText = "未采集信息"
	case "1":
		_self.AccountTypeText = "正常"
	case "2":
		_self.AccountTypeText = "已过期"
	case "3":
		_self.AccountTypeText = "已禁用"
	}

	switch _self.AccountStatus {
	case "-1":
		_self.AccountStatusText = "采集失败"
	case "0":
		_self.AccountStatusText = "特权帐号"
	case "1":
		_self.AccountStatusText = "管理帐号"
	case "2":
		_self.AccountStatusText = "普通帐号"
	}

}

// AssetAccountExtend 从帐号扩展表
type AssetAccountExtend struct {
	ModelId
	Uid           string      `gorm:"size:32;comment:从帐号Uid"  json:"uid"`
	HomePath      string      `gorm:"size:256;comment:home路径"  json:"homePath"`
	Shell         string      `gorm:"size:256;comment:用户使用的shell"  json:"shell"`
	LastUpPwdTime c_type.Time `gorm:"type:date;comment:最后修改密码时间"  json:"lastUpPwdTime"`
	ExpirePwdTime c_type.Time `gorm:"type:date;comment:密码过期天数"  json:"expirePwdTimeAt"`
	InactiveAt    c_type.Time `gorm:"type:date;comment:过期后缓冲期"  json:"inactiveAtAt"`
	ExpireAt      c_type.Time `gorm:"type:date;comment:帐号失效时间"  json:"expireAtAt"`
	Remark        string      `gorm:"size:256;comment:备注"  json:"remark"`
	CollectTime   c_type.Time `gorm:"size:512;comment:最后采集时间"  json:"collectTime"`
	CollectMsg    string      `gorm:"size:512;comment:采集结果信息"  json:"collectMsg"`
	//RawPasswd     string `gorm:"size:4096;comment:原始记录-/etc/passwd"  json:"rawPasswd"`
	//RawShadow     string `gorm:"size:4096;comment:原始记录-/etc/shadow"  json:"rawShadow"`
	//RawGroup      string `gorm:"size:4096;comment:原始记录-/etc/group"  json:"rawGroup"`
	ModelControl
	ModelDelete
}

func (AssetAccountExtend) TableName() string {
	return table_name.AssetAccountExtend
}
