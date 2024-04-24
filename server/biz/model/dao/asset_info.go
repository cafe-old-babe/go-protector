package dao

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/entity"
	"gorm.io/gorm"
)

var AssetBasic assetBasic

type assetBasic struct {
}

func (_self assetBasic) CheckSave(db *gorm.DB, model entity.AssetBasic) (err error) {
	if err = binding.Validator.ValidateStruct(model); err != nil {
		return err
	}
	var count int64
	err = db.Model(&model).Scopes(
		func(db *gorm.DB) *gorm.DB {
			if model.ID >= 0 {
				db = db.Where("id <> ?", model.ID)
			}
			db = db.Where("(asset_name = ? or ip = ?)", model.AssetName, model.IP)
			return db
		},
	).Limit(-1).Offset(-1).Count(&count).Error
	if err == nil && count > 0 {
		err = errors.New("资产与IP不可重复")
	}
	return
}
