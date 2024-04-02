package service

import (
	"errors"
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/scope"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"gorm.io/gorm"
)

type AssetNetwork struct {
	base.Service
}

func (_self *AssetNetwork) Page(req *dto.AssetNetworkPageReq) (res *base.Result) {
	if req == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var slice []entity.AssetNetwork
	var count int64
	if err := _self.DB.Scopes(
		scope.Paginate(req),
		scope.Like("an_name", req.AnName),
		scope.Like("an_ip", req.AnIp),
	).Find(&slice).Limit(-1).Offset(-1).
		Count(&count).Error; err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	res = base.ResultPage(slice, req, count)
	return
}

func (_self *AssetNetwork) Check(model *entity.AssetNetwork) (err error) {

	var count int64
	if err = _self.DB.Model(model).Scopes(
		func(db *gorm.DB) *gorm.DB {
			if model.ID > 0 {
				db = db.Where("id <> ?", model.ID)
			}
			return db.Where("(an_name = ? or an_ip = ?)", model.AnName, model.AnIp)
		}).Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		err = errors.New("名称和Ip不能重复")
	}
	return

}
