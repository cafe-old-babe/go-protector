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
			condition.Eq(table_name.SsoSession+".sso_session_id ", req.SsoSessionId),
			condition.Eq(table_name.SsoSession+".asset_ip ", req.AssetIp),
			condition.Eq(table_name.SsoSession+".asset_name ", req.AssetName),
			condition.Like("cmd ", req.Cmd),
		).Order(table_name.SsoSession + ".created_at").
			Order(clause.OrderByColumn{Column: clause.Column{Name: "sort"}, Desc: false}).
			Joins("SsoSession").Find(&slice))
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	res = base.ResultPage(slice, req, count)
	return

}
