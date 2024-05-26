package service

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/database/condition"
)

type AssetAuth struct {
	base.Service
}

// Page 分页查询
func (_self *AssetAuth) Page(req *dto.AssetAuthPageReq) (res *base.Result) {

	var slice []entity.AssetAuth
	count, err := _self.Count(
		_self.DB.Scopes(
			condition.Paginate(req),
			condition.Like("asset_ip", req.AssetIp),
			condition.Like("asset_acc", req.AssetAcc),
			condition.Like("user_acc", req.UserAcc),
		).Find(&slice))
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	res = base.ResultPage(slice, req, count)
	return
}

// SaveCheck 保存前更新
func (_self *AssetAuth) SaveCheck(data *entity.AssetAuth) (err error) {
	if err = binding.Validator.ValidateStruct(data); err != nil {
		return err
	}

	count, err := _self.Count(
		_self.DB.Model(data).
			Where("asset_id = ? and asset_acc_id = ? and user_id",
				data.AssetId, data.AssetAccId, data.UserId),
	)
	if err != nil {
		return err
	}
	if count > 0 {
		err = errors.New("授权重复")
	}

	return
}
