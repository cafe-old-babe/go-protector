package service

import (
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/scope"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
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

// Insert 字典类型新增
func (_self *SysDictType) Insert(model *entity.SysDictType) *dto.Result {
	if err := _self.DB.Create(model).Error; err != nil {
		return dto.ResultFailureErr(err)
	}

	return dto.ResultSuccess(model, "创建成功")
}

// Update 字典类型更新
func (_self *SysDictType) Update(model *entity.SysDictType) *dto.Result {
	if err := _self.DB.Save(model).Error; err != nil {
		return dto.ResultFailureErr(err)
	}
	return dto.ResultSuccess(model, "更新成功")
}

// Delete 字典类型删除
func (_self *SysDictType) Delete(req *dto.IdsReq) *dto.Result {
	if req == nil || len(req.Ids) <= 0 {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
	}
	result := _self.DB.Delete(&entity.SysDictType{}, req.Ids)
	if result.Error != nil {
		return dto.ResultFailureErr(result.Error)
	}
	if result.RowsAffected <= 0 {
		_self.Logger.Error("删除失败,无删除记录")
		return dto.ResultFailureMsg("删除失败")
	}
	return dto.ResultSuccessMsg("删除成功")
}
