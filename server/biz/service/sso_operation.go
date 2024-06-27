package service

import (
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
	"gorm.io/gorm/clause"
)

type SsoOperation struct {
	base.Service
}

func (_self SsoOperation) Page(req *dto.SsoOperationPageReq) (res *base.Result) {
	if req == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var slice []entity.SsoOperation
	count, err := _self.Count(
		_self.GetDB().Scopes(
			condition.Paginate(req),
			condition.Like("SsoSession.asset_acc ", req.AssetAcc),
			condition.Like("SsoSession.asset_ip ", req.AssetIp),
			condition.Like("SsoSession.asset_name ", req.AssetName),
			condition.Like("cmd ", req.Cmd),
		).Order("SsoSession.created_at desc").
			Order(clause.OrderByColumn{Column: clause.Column{Name: table_name.SsoOperation + ".sort"}, Desc: false}).
			Joins("SsoSession").Find(&slice))
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	res = base.ResultPage(slice, req, count)
	return

}
