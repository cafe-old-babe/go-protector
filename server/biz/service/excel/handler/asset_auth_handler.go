package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/dao"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"gorm.io/gorm"
)

type AssetAuthExcelHandler[T entity.AssetAuth] struct {
	ErrData []entity.AssetAuth
	base.Service
}

func (_self *AssetAuthExcelHandler[T]) ReadRow(row *entity.AssetAuth) (err error) {
	if row == nil {
		return
	}
	if err = fillAuth(_self.DB, row); err != nil {
		return
	}
	if err = binding.Validator.ValidateStruct(row); err != nil {
		return
	}
	if err = _self.DB.Create(row).Error; err != nil {
		err = errors.New(fmt.Sprintf("新增授权失败: %s", err.Error()))
	}

	return
}

func fillAuth(db *gorm.DB, auth *entity.AssetAuth) (err error) {
	assetBasic, err := dao.AssetBasic.FindAssetBasicByDTO(db, dto.FindAssetDTO{
		AssetName: auth.AssetName,
		IP:        auth.AssetIp,
	})
	if err != nil {
		return
	}
	auth.AssetId = assetBasic.ID

	assetAccount, err := dao.AssetAccount.FindAssetAccountByDTO(db, dto.FindAssetAccountDTO{
		AssetId: auth.AssetId,
		Account: auth.AssetAcc,
	})
	if err != nil {
		return
	}
	auth.AssetAccId = assetAccount.ID

	sysUser, err := dao.SysUser.FindUserByDTO(db, dto.FindUserDTO{
		LoginName: auth.UserAcc,
	})
	if err != nil {
		return
	}
	auth.UserId = sysUser.ID
	return
}

func (_self *AssetAuthExcelHandler[T]) ReadDone() {

}

func (_self *AssetAuthExcelHandler[T]) NewRow() *entity.AssetAuth {
	return &entity.AssetAuth{}
}

func (_self *AssetAuthExcelHandler[T]) AppendErrData(row *entity.AssetAuth) {
	_self.ErrData = append(_self.ErrData, *row)
}
