package dao

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
	"gorm.io/gorm"
	"reflect"
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

// FindAssetBasicByDTO 根据DTO查询资产信息
func (_self assetBasic) FindAssetBasicByDTO(db *gorm.DB, param dto.FindAssetDTO) (model entity.AssetBasic, err error) {
	if reflect.ValueOf(param).IsZero() {
		err = c_error.ErrParamInvalid
		return
	}

	err = db.Scopes(
		condition.Eq("asset_name", param.AssetName),
		condition.Eq("ip", param.IP),
		condition.Eq("id", param.ID),
	).First(&model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("未查询到资产信息")
	}
	return
}
