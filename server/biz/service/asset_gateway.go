package service

import (
	"errors"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
	"gorm.io/gorm"
	"sync"
)

type AssetGateway struct {
	base.Service
}

var gateWayDTOCache sync.Map

func (_self *AssetGateway) Page(req *dto.AssetGatewayPageReq) (res *base.Result) {
	if req == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var slice []entity.AssetGateway
	var count int64
	if err := _self.DB.Scopes(
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
	if err = _self.DB.Model(model).Scopes(
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

func (_self *AssetGateway) GetGatewayDTOById(id uint64) (gatewayDTO *dto.GatewayDTO, err error) {
	if id <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	if value, ok := gateWayDTOCache.Load(id); ok {
		gatewayDTO = value.(*dto.GatewayDTO)
		return
	}

	return
}
