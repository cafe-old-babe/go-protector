package dao

import (
	"errors"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/consts/table_name"
	"gorm.io/gorm"
)

var AssetAuth assetAuth

type assetAuth struct {
}

func (_self assetAuth) GenerateSubQueryForAssetId(db *gorm.DB, userId uint64) *gorm.DB {
	if userId <= 0 {
		return db
	}
	tx := db.Session(&gorm.Session{NewDB: true})

	return tx.Select("asset_id").
		Table(table_name.AssetAuth).
		Where(" user_id = ? ", userId).
		Where("deleted_at is null").
		Where(" sysdate() between IFNULL(start_date,sysdate()) and IFNULL(end_date,sysdate()) ")
}

func (_self assetAuth) GenerateSubQueryForAssetAccId(db *gorm.DB, userId uint64, isBetween bool) *gorm.DB {
	if userId <= 0 {
		return db
	}
	tx := db.Session(&gorm.Session{NewDB: true})
	tx = tx.Select("asset_acc_id").
		Table(table_name.AssetAuth).
		Where("user_id = ? ", userId).
		Where("deleted_at is null")
	if isBetween {
		tx = tx.Where(" sysdate() between IFNULL(start_date,sysdate()) and IFNULL(end_date,sysdate()) ")
	}
	return tx
}

func (_self assetAuth) FindById(db *gorm.DB, id uint64) (model entity.AssetAuth, err error) {
	if err = db.First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("未查询到授权信息")
		}
	}

	return
}
