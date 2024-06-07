package dao

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
	"gorm.io/gorm"
)

var AssetAccount assetAccount

type assetAccount struct {
}

func (_self assetAccount) CheckSave(db *gorm.DB, model *entity.AssetAccount) (err error) {
	if model.AccountType == "0" && model.ID <= 0 {
		// 特权帐号特殊处理,保证特权帐号只有一个
		err = db.Model(&entity.AssetAccount{}).Scopes(func(db *gorm.DB) *gorm.DB {
			if model.AssetId > 0 {
				db = db.Where("asset_id = ? ", model.AssetId)
			}
			return db.Where("account_type = '0'")
		}).Limit(1).Pluck("id", &model.ID).Error
		if err != nil {
			return err
		}
	}
	if err = binding.Validator.ValidateStruct(model); err != nil {
		return err
	}
	var count int64
	var errMsg string
	// 特权账号
	if model.AccountType == "0" {
		db = db.Model(model).Scopes(
			func(db *gorm.DB) *gorm.DB {
				db = db.Where("account_type = ?", "0")
				db = db.Where("asset_id = ?", model.AssetId)
				if model.ID > 0 {
					db = db.Where("id <> ?", model.ID)
				}
				return db
			},
		)
		errMsg = "保存失败,该资产下已存在特权帐号,不可再次添加"
	} else {

		db = db.Model(model).Scopes(
			func(db *gorm.DB) *gorm.DB {
				db = db.Where("account = ?", model.Account)
				db = db.Where("asset_id = ?", model.AssetId)
				if model.ID > 0 {
					db = db.Where("id <> ?", model.ID)
				}
				return db
			},
		)
		errMsg = "保存失败,资产下已存在[" + model.Account + "],不可重复添加"
	}
	err = db.Limit(-1).Offset(-1).Count(&count).Error

	if err == nil && count > 0 {
		err = errors.New(errMsg)
	}

	return
}

func (_self assetAccount) DeleteByAssetId(db *gorm.DB, ids []uint64) error {
	if len(ids) <= 0 {
		return c_error.ErrParamInvalid
	}
	return db.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Delete(&entity.AssetAccount{}, ids).Error; err != nil {
			return err
		}
		if err = tx.Delete(&entity.AssetAccountExtend{}, ids).Error; err != nil {
			return err
		}
		return
	})

}

// FindAssetAccountByDTO 查询从账号
func (_self assetAccount) FindAssetAccountByDTO(db *gorm.DB, param dto.FindAssetAccountDTO) (
	model entity.AssetAccount, err error) {

	if err = binding.Validator.ValidateStruct(param); err != nil {
		return
	}
	if err = db.Scopes(
		condition.Eq("asset_id", param.AssetId),
		condition.Eq("account", param.Account),
	).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("未查询到从账号信息")
		}
	}

	return
}
