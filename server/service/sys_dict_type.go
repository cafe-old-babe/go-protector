package service

import (
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/scope"
	"go-protector/server/dao"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"gorm.io/gorm"
)

type SysDictType struct {
	base.Service
}

// Page 字典类型分页查询
func (_self *SysDictType) Page(req *dto.DictTypePageReq) (res *dto.Result) {
	var dictType entity.SysDictType
	var list []entity.SysDictType
	var count int64
	if err := _self.DB.Model(&dictType).
		Scopes(
			scope.Paginate(req),
			scope.Like("type_code", req.TypeCode),
			scope.Like("type_name", req.TypeName),
		).Find(&list).                      // 查询数据
		Limit(-1).Offset(-1).Count(&count). // 查询总数
		Error; err != nil {
		return dto.ResultFailureErr(err)
	}
	return dto.ResultPage(list, req, count)
}

// Save 字典类型保存
func (_self *SysDictType) Save(model *entity.SysDictType) *dto.Result {
	// todo 校验
	if err := _self.DB.Save(model).Error; err != nil {
		return dto.ResultFailureErr(err)
	}
	return dto.ResultSuccess(model, "创建成功")
}

// Update 字典类型更新
func (_self *SysDictType) Update(model *entity.SysDictType) *dto.Result {
	// todo 校验
	if err := _self.DB.Save(model).Error; err != nil {
		return dto.ResultFailureErr(err)
	}
	return dto.ResultSuccess(model, "更新成功")
}

// Delete 字典类型删除
func (_self *SysDictType) Delete(req *dto.IdsReq) *dto.Result {

	if req == nil || len(req.GetIds()) <= 0 {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
	}
	ids := req.GetIds()
	//todo 同时删除data
	err := _self.DB.Transaction(func(tx *gorm.DB) error {
		var dictTypeSlice []entity.SysDictType
		if err := tx.Find(&dictTypeSlice, ids).Error; err != nil {
			return err
		}
		if len(dictTypeSlice) <= 0 || len(dictTypeSlice) != len(ids) {
			return c_error.ErrDeleteFailure
		}
		var idSlice []uint64
		var typeCodeSlice []string
		for _, elem := range dictTypeSlice {
			idSlice = append(idSlice, elem.ID)
			typeCodeSlice = append(typeCodeSlice, elem.TypeCode)
		}
		if err := dao.SysDict.DeleteTypeByIds(tx, ids); err != nil {
			return err
		}
		return dao.SysDict.DeleteDataByTypeCode(tx, typeCodeSlice)
	})

	if err != nil {
		return dto.ResultFailureErr(err)
	}
	return dto.ResultSuccessMsg("删除成功")
}
