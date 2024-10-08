package service

import (
	"errors"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
	"gorm.io/gorm"
)

type AssetGateway struct {
	base.Service
}

func (_self *AssetGateway) Page(req *dto.AssetGatewayPageReq) (res *base.Result) {
	if req == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var slice []entity.AssetGateway
	var count int64
	if err := _self.GetDB().Scopes(
		condition.Paginate(req),
		condition.Like("ag_name", req.AgName),
		condition.Like("ag_ip", req.AgIp),
	).Find(&slice).Limit(-1).Offset(-1).
		Count(&count).Error; err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	res = base.ResultPage(slice, req, count)
	return
}

func (_self *AssetGateway) Check(model *entity.AssetGateway) (err error) {

	var count int64
	if err = _self.GetDB().Model(model).Scopes(
		func(db *gorm.DB) *gorm.DB {
			if model.ID > 0 {
				db = db.Where("id <> ?", model.ID)
			}
			return db.Where("(ag_name = ? or ag_ip = ?)", model.AgName, model.AgIp)
		}).Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		err = errors.New("名称和Ip不能重复")
	}
	return

}

func (_self *AssetGateway) List() (res *base.Result) {
	var slice []entity.AssetGateway

	if err := _self.GetDB().Find(&slice).Error; err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	for i := range slice {
		slice[i].AgPassword = ""
	}
	res = base.ResultSuccess(slice)
	return
}
