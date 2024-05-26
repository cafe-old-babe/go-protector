package entity

import (
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/database/condition"
	"gorm.io/gorm"
	"reflect"
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

// UpdateRedundancy 更新冗余数据
func (_self *AssetAuth) UpdateRedundancy(db *gorm.DB, data interface{}) (err error) {
	if err = db.Error; err != nil {
		return
	}
	if db == nil || data == nil {
		err = c_error.ErrParamInvalid
		return
	}

	valueOf := reflect.ValueOf(data)
	indirect := reflect.Indirect(valueOf)
	valueOfType := indirect.Type()
	// 获取真实的kind类型
	if valueOfType.Kind() != reflect.Struct {
		err = c_error.ErrParamInvalid
		return
	}
	tx := db.Session(&gorm.Session{NewDB: true})

	if indirect.CanConvert(TypeSysUser) {
		if !db.Statement.Changed("LoginName") {
			return
		}
		convert := indirect.Convert(TypeSysUser)
		id := convert.FieldByName("ID")
		if id.IsZero() {
			return
		}
		userAcc := convert.FieldByName("LoginName")
		_self.UserId = id.Uint()
		_self.UserAcc = userAcc.String()
		tx = tx.Select("user_acc").Where("user_id = ? ", _self.UserId)
	} else if indirect.CanConvert(TypeAssetBasic) {
		if !db.Statement.Changed("IP", "AssetName") {
			return
		}
		convert := indirect.Convert(TypeAssetBasic)
		id := convert.FieldByName("ID")
		if id.IsZero() {
			return
		}
		ip := convert.FieldByName("IP")
		assetName := convert.FieldByName("AssetName")
		_self.AssetId = id.Uint()
		_self.AssetIp = ip.String()
		_self.AssetName = assetName.String()

		tx = tx.Select("ip", "asset_name").
			Where("asset_id = ? ", _self.AssetId)

	} else if indirect.CanConvert(TypeAssetAccount) {
		if !db.Statement.Changed("Account") {
			return
		}

		convert := indirect.Convert(TypeAssetAccount)
		id := convert.FieldByName("ID")
		if id.IsZero() {
			return
		}
		assetAcc := convert.FieldByName("Account")
		_self.AssetAccId = id.Uint()
		_self.AssetAcc = assetAcc.String()
		tx = tx.Select("asset_acc").
			Where("asset_acc_id = ? ", _self.AssetAccId)

	} else {
		err = c_error.ErrParamInvalid
		return
	}
	// https://gorm.io/zh_CN/docs/update.html#%E4%B8%8D%E4%BD%BF%E7%94%A8-Hook-%E5%92%8C%E6%97%B6%E9%97%B4%E8%BF%BD%E8%B8%AA
	if err = tx.UpdateColumns(_self).Error; err != nil {
		return
	}

	return
}

// DeleteRedundancy 删除冗余数据
func (_self *AssetAuth) DeleteRedundancy(db *gorm.DB, ids []uint64, dateType reflect.Type) (err error) {
	if db == nil || ids == nil || len(ids) <= 0 || dateType == nil {
		return
	}
	tx := db.Session(&gorm.Session{NewDB: true})
	var scope func(*gorm.DB) *gorm.DB
	switch dateType {
	case TypeSysUser:
		scope = condition.In[uint64]("user_id", ids)
		//tx = tx.Where("user_id in (?)")
	case TypeAssetBasic:
		scope = condition.In[uint64]("asset_id", ids)
		//tx = tx.Where("asset_id in (?)", ids)

	case TypeAssetAccount:
		scope = condition.In[uint64]("asset_acc_id", ids)
		//tx = tx.Where("asset_acc_id in (?)", ids)
	default:
		err = c_error.ErrParamInvalid
	}
	return tx.Scopes(scope).Delete(_self).Error
}
