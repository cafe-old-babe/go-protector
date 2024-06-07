package service

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service/excel/handler"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
	"go-protector/server/internal/database/condition"
	"go-protector/server/internal/utils/excel"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
	"strings"
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
		_self.DB.Model(data).Scopes(func(db *gorm.DB) *gorm.DB {
			if data.ID > 0 {
				db = db.Where("id <> ?", data.ID)
			}
			return db
		}).Where("asset_id = ? and asset_acc_id = ? and user_id = ?",
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

// Import 导入
func (_self *AssetAuth) Import(file *multipart.FileHeader) (err error) {
	_self.Logger.Debug("fileName: %s", file.Filename)

	open, err := file.Open()
	if err != nil {
		return

	}
	defer open.Close()
	var excelHandler handler.AssetAuthExcelHandler[entity.AssetAuth]
	_self.MakeService(&excelHandler)
	if err = excel.ReadExcelFirstSheet[*entity.AssetAuth](open, &excelHandler); err != nil {
		return
	}

	if len(excelHandler.ErrData) <= 0 {
		c_result.Success(_self.Context, nil, "导入成功")
		return
	}
	ext := filepath.Ext(file.Filename)
	fileName := strings.ReplaceAll(file.Filename, ext, "_err"+ext)
	resFile, err := excel.GenerateExcel(excelHandler.ErrData)
	if err != nil {
		return
	}
	defer resFile.Close()
	err = excel.Export(_self.Context, resFile, fileName)
	return
}

func (_self *AssetAuth) ExportData(req *dto.AssetAuthPageReq) (err error) {

	var slice []entity.AssetAuth
	err = _self.DB.Scopes(
		condition.Paginate(req),
		condition.Like("asset_ip", req.AssetIp),
		condition.Like("asset_acc", req.AssetAcc),
		condition.Like("user_acc", req.UserAcc),
	).Find(&slice).Error
	if err != nil {
		return
	}
	if len(slice) <= 0 {
		err = errors.New("未查询到数据")
		return
	}

	file, err := excel.GenerateExcel(slice, "错误消息")
	if err != nil {
		return
	}
	err = excel.Export(_self.Context, file, "授权数据.xlsx")
	return
}

func (_self *AssetAuth) ExportTemplate() (err error) {

	file, err := excel.GenerateExcel(entity.AssetAuth{}, "错误消息")
	if err != nil {
		return
	}
	err = excel.Export(_self.Context, file, "授权导入模板.xlsx")
	return
}
