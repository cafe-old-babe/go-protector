package service

import (
	"errors"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service/iface"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
)

type ApproveCmd struct {
	base.Service
}

func init() {
	iface.RegisterApproveCmdService(&ApproveCmd{})
}

func (_self *ApproveCmd) Page(req *dto.ApproveCmdPageReq) (res *base.Result) {
	var slice []entity.ApproveCmd
	count, err := _self.Count(
		_self.GetDB().Scopes(
			condition.Paginate(req),
			condition.Eq("asset_id", req.AssetId),
			condition.Like("cmd", req.Cmd),
		).Find(&slice))
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	res = base.ResultPage(slice, req, count)
	return
}

func (_self *ApproveCmd) SaveCheck(req entity.ApproveCmd) (err error) {
	if len(req.Cmd) <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	var count int64
	err = _self.GetDB().Table(table_name.ApproveCmd).
		Where("asset_id = ? and cmd = ?", req.AssetId, req.Cmd).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		err = errors.New("指令重复")
	}
	return
}

func (_self *ApproveCmd) GetApproveCmdSliceByAssetId(assetId uint64) (cmdSlice []string, err error) {
	if assetId <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	err = _self.GetDB().Table(table_name.ApproveCmd).Where("asset_id = ?", assetId).Pluck("cmd", &cmdSlice).Error
	return
}
