package entity

import (
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/database/condition"
	"go-protector/server/internal/utils/excel"
	"gorm.io/gorm"
	"reflect"
)

// AssetAuth 授权表
type AssetAuth struct {
	ModelId
	AssetId      uint64       `gorm:"comment:资产id;uniqueIndex:auth_unique_i"  json:"assetId" binding:"required"`
	AssetName    string       `gorm:"size:32;comment:资产名称" json:"assetName" binding:"required"  excel:"title:资产名称" `
	AssetIp      string       `gorm:"size:32;comment:资产IP" json:"assetIp" binding:"required" excel:"title:资产IP"`
	AssetAccId   uint64       `gorm:"comment:资产从帐号ID;uniqueIndex:auth_unique_i"  json:"assetAccId" binding:"required"`
	AssetAcc     string       `gorm:"size:32;comment:资产从帐号" json:"assetAcc" binding:"required" excel:"title:资产从帐号"`
	UserId       uint64       `gorm:"comment:主帐号ID;uniqueIndex:auth_unique_i"  json:"userId" binding:"required"`
	UserAcc      string       `gorm:"size:32;comment:主帐号"  json:"userAcc" binding:"required" excel:"title:主帐号"`
	StartDate    c_type.Time  `gorm:"type:date;comment:授权开始时间"  json:"startDate" excel:"title:授权开始时间"`
	EndDate      c_type.Time  `gorm:"type:date;comment:授权结束时间"  json:"endDate" excel:"title:授权结束时间"`
	AssetAccount AssetAccount `gorm:"foreignKey:AssetAccId;references:ID"  json:"assetAccount,omitempty" binding:"-"`
	ModelControl
	ModelDelete
	excel.StdRow `gorm:"-" json:"-"`
}

func (*AssetAuth) TableName() string {
	return table_name.AssetAuth
}

// UpdateRedundancy 更新冗余数据
// 7-4	【实战】使用钩子函数更新反范式的冗余数据-掌握通过反射获取interface的数据类型与值
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
	return db.Transaction(func(tx *gorm.DB) (err error) {

		if indirect.CanConvert(TypeSysUser) {

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

			var selectSlice []string
			convert := indirect.Convert(TypeAssetBasic)
			id := convert.FieldByName("ID")
			if id.IsZero() {
				return
			}

			if ip := convert.FieldByName("IP"); !ip.IsZero() {
				_self.AssetIp = ip.String()
				selectSlice = append(selectSlice, "asset_ip")
			}
			if assetName := convert.FieldByName("AssetName"); !assetName.IsZero() {
				_self.AssetName = assetName.String()
				selectSlice = append(selectSlice, "asset_name")
			}
			if len(selectSlice) <= 0 {
				return
			}
			tx = tx.Select(selectSlice).
				Where("asset_id = ? ", id.Interface())

		} else if indirect.CanConvert(TypeAssetAccount) {

			convert := indirect.Convert(TypeAssetAccount)
			id := convert.FieldByName("ID")
			if id.IsZero() {
				return
			}
			assetAcc := convert.FieldByName("Account")
			if assetAcc.IsZero() {
				return
			}
			_self.AssetAcc = assetAcc.String()
			tx = tx.Select("asset_acc").
				Where("asset_acc_id = ? ", id.Interface())

		} else {
			err = c_error.ErrParamInvalid
			return
		}
		// https://gorm.io/zh_CN/docs/update.html#%E4%B8%8D%E4%BD%BF%E7%94%A8-Hook-%E5%92%8C%E6%97%B6%E9%97%B4%E8%BF%BD%E8%B8%AA

		err = tx.Model(_self).UpdateColumns(_self).Error
		return
	})
}

// DeleteRedundancy 删除冗余数据
// 7-5	【实战】删除反范式的冗余数据-掌握如何解决在in与not in数量过多的场景下导致失败的问题
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
	if err != nil {
		return
	}
	return tx.Scopes(scope).Delete(_self).Error
}

func (_self *AssetAuth) AfterFind(db *gorm.DB) error {
	if err := _self.AssetAccount.AfterFind(db); err != nil {
		return err
	}
	_self.AssetAccount.Password = ""
	return nil
}
