package service

import (
	"errors"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
	"gorm.io/gorm"
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

func (_self SsoOperation) PageBySsoId(req *dto.SsoOperationPageReq) (res *base.Result) {
	if req == nil || req.SsoSessionId <= 0 {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var err error
	var slice []dto.SsoOperationBySsoIdPage

	var ssoSession entity.SsoSession
	if err = _self.GetDB().First(&ssoSession, req.SsoSessionId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res = base.ResultFailureErr(c_error.ErrIllegalAccess)
		} else {
			res = base.ResultFailureErr(err)
		}
		return
	}
	if !ssoSession.ConnectAt.Valid || ssoSession.Status != consts.SessionClose {
		res = base.ResultFailureErr(errors.New("会话未开始或未结束"))
		return
	}

	var count int64
	count, err = _self.Count(_self.GetDB().Table("sso_operation").
		Scopes(
			condition.Paginate(req),
			condition.Eq("sso_session_id", req.SsoSessionId),
			condition.Like("cmd", req.Cmd),
		).Order("sort").
		Select("id", "sso_session_id", "cmd", "cmd_start_at").
		Find(&slice))

	if err != nil {
		res = base.ResultFailureErr(err)
	}
	for i := range slice {
		slice[i].TimeStamp = slice[i].CmdStartAt.Time.Sub(ssoSession.ConnectAt.Time).Seconds()
	}
	res = base.ResultPage(slice, req, count)
	return
}
