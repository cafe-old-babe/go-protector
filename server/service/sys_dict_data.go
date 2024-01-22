package service

import (
	"go-protector/server/core/base"
	"go-protector/server/core/consts/table_name"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/scope"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"go-protector/server/models/vo"
	"gorm.io/gorm"
)

type SysDictData struct {
	base.Service
}

// Page 字典数据分页查询
func (_self *SysDictData) Page(req *dto.DictDataPageReq) *dto.Result {
	var dictData entity.SysDictData
	var list []vo.DictDataPage
	var count int64
	if err := _self.DB.Model(&dictData).
		Select([]string{
			table_name.SysDictData + ".id",
			table_name.SysDictType + ".type_name",
			table_name.SysDictType + ".type_code",
			table_name.SysDictData + ".data_name",
			table_name.SysDictData + ".data_code",
			table_name.SysDictData + ".type_status",
			`case  when ` + table_name.SysDictData + `.type_status  = 0 then '正常' when ` +
				table_name.SysDictData + `.type_status  = 1 then '停用' end as type_status_text`,
		}).Scopes(
		scope.Paginate(req),
		scope.Like("data_code", req.DataCode),
		scope.Like("data_name", req.DataName),
		func(db *gorm.DB) *gorm.DB {
			if len(req.TypeCode) <= 0 {
				return db
			}
			return db.Where("type_code = ?", req.TypeCode)
		}).Order("sort").
		Joins(`left join ` + table_name.SysDictType + `on ` +
							table_name.SysDictType + `.type_code = ` + table_name.SysDictData + `.type_code`).
		Find(&list).                        // 查询数据
		Limit(-1).Offset(-1).Count(&count). // 查询总数
		Error; err != nil {
		return dto.ResultFailureErr(err)
	}
	return dto.ResultPage(list, req, count)
}

// Insert 字典数据新增
func (_self *SysDictData) Insert(model *entity.SysDictData) *dto.Result {
	if err := _self.DB.Create(model).Error; err != nil {
		return dto.ResultFailureErr(err)
	}

	return dto.ResultSuccess(model, "创建成功")
}

// Update 字典数据更新
func (_self *SysDictData) Update(model *entity.SysDictData) *dto.Result {
	if err := _self.DB.Save(model).Error; err != nil {
		return dto.ResultFailureErr(err)
	}
	return dto.ResultSuccess(model, "更新成功")
}

// UpdateStatus 更新状态
func (_self *SysDictData) UpdateStatus(req *dto.DictDataUpdateStatusReq) *dto.Result {

	result := _self.DB.Table(table_name.SysDictData, req.ID).
		Update("type_status", req.Status)

	if result.Error != nil {
		return dto.ResultFailureErr(result.Error)
	}
	if result.RowsAffected <= 0 {
		_self.Logger.Error("更新失败,无更新记录")
		return dto.ResultFailureMsg("更新失败")
	}
	return dto.ResultSuccessMsg("更新成功")

}

// Delete 字典数据删除
func (_self *SysDictData) Delete(req *dto.IdsReq) *dto.Result {
	if req == nil || len(req.Ids) <= 0 {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
	}
	result := _self.DB.Delete(&entity.SysDictData{}, req.Ids)
	if result.Error != nil {
		return dto.ResultFailureErr(result.Error)
	}
	if result.RowsAffected <= 0 {
		_self.Logger.Error("删除失败,无删除记录")
		return dto.ResultFailureMsg("删除失败")
	}
	return dto.ResultSuccessMsg("删除成功")

}
