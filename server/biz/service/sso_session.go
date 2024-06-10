package service

import (
	"errors"
	"go-protector/server/biz/model/dao"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
)

type SsoSession struct {
	base.Service
}

// CreateSession 创建会话
func (_self *SsoSession) CreateSession(authId uint64) (res *base.Result) {
	if authId <= 0 {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	// 查询授权信息
	var auth entity.AssetAuth
	var err error
	if auth, err = dao.AssetAuth.FindById(_self.DB, authId); err != nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	if auth.StartDate.Valid && auth.EndDate.Valid {
		now := c_type.NowTime()
		if now.Time.Before(auth.StartDate.Time) || now.Time.After(auth.EndDate.Time) {
			res = base.ResultFailureErr(errors.New("授权无效或已过期"))
			return
		}
	}
	// 校验授权信息
	if auth.UserId != current.GetUserId(_self.Context) {
		res = base.ResultFailureErr(errors.New("授权无效或已过期"))
		return
	}

	assetBasic, err := dao.AssetBasic.FindAssetBasicByDTO(_self.DB, dto.FindAssetDTO{
		ID: auth.AssetId,
	})
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	assetAccount, err := dao.AssetAccount.FindAssetAccountByDTO(_self.DB, dto.FindAssetAccountDTO{
		AssetId: auth.AssetId,
		Account: auth.AssetAcc,
	})
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}

	// 创建
	var session entity.SsoSession
	session.AuthId = auth.ID

	session.AssetId = auth.AssetId
	session.AssetName = auth.AssetName
	session.AssetIp = auth.AssetIp
	session.AssetPort = assetBasic.Port
	session.AssetGatewayId = assetBasic.AssetGatewayId

	session.AssetAccId = auth.AssetAccId
	session.AssetAcc = auth.AssetAcc
	session.AssetAccPwd = assetAccount.Password

	session.UserId = auth.UserId
	session.UserAcc = auth.UserAcc

	session.Status = "0"

	if err = _self.DB.Create(&session).Error; err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	return base.ResultSuccess(map[string]uint64{"id": session.ID}, "创建成功")
}

func (_self *SsoSession) ConnectSession(req *dto.ConnectSessionReq) (err error) {

	return
}
