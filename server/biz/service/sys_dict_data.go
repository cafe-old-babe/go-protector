package service

import (
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/model/vo"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
	"gorm.io/gorm"
)

type SysDictData struct {
	base.Service
}

// Page 字典数据分页查询
func (_self *SysDictData) Page(req *dto.DictDataPageReq) *base.Result {
	if len(req.TypeCode) <= 0 {
		return base.ResultFailureMsg("请选择字典类型")
	}

	var dictData entity.SysDictData
	var list []vo.DictDataPage
	var count int64
	if err := _self.GetDB().Model(&dictData).
		Select([]string{
			table_name.SysDictData + ".id",
			table_name.SysDictType + ".type_name",
			table_name.SysDictType + ".type_code",
			table_name.SysDictData + ".data_name",
			table_name.SysDictData + ".data_code",
			table_name.SysDictData + ".status",
			`case  when ` + table_name.SysDictData + `.status  = 0 then '正常' when ` +
				table_name.SysDictData + `.status  = 1 then '停用' end as status_text`,
		}).Scopes(
		condition.Paginate(req),
		condition.Like("data_code", req.DataCode),
		condition.Like("data_name", req.DataName),
		func(db *gorm.DB) *gorm.DB {
			if len(req.TypeCode) <= 0 {
				return db
			}
			return db.Where(table_name.SysDictData+".type_code = ?", req.TypeCode)
		}).Order("sort").
		Joins(`left join ` + table_name.SysDictType + ` on ` +
							table_name.SysDictType + `.type_code = ` + table_name.SysDictData + `.type_code`).
		Find(&list).                        // 查询数据
		Limit(-1).Offset(-1).Count(&count). // 查询总数
		Error; err != nil {
		return base.ResultFailureErr(err)
	}
	return base.ResultPage(list, req, count)
}

// Save 字典数据新增
func (_self *SysDictData) Save(model *entity.SysDictData) *base.Result {
	// todo 校验
	if err := _self.GetDB().Save(model).Error; err != nil {
		return base.ResultFailureErr(err)
	}

	return base.ResultSuccess(model, "创建成功")
}

// UpdateStatus 更新状态
func (_self *SysDictData) UpdateStatus(req *dto.DictDataUpdateStatusReq) *base.Result {

	result := _self.GetDB().Model(&entity.SysDictData{}).
		Where("id = ?", req.ID).
		Update("status", req.Status)

	if result.Error != nil {
		return base.ResultFailureErr(result.Error)
	}
	if result.RowsAffected <= 0 {
		_self.GetLogger().Error("更新失败,无更新记录")
		return base.ResultFailureMsg("更新失败")
	}
	return base.ResultSuccessMsg("更新成功")

}

// Delete 字典数据删除
func (_self *SysDictData) Delete(req *base.IdsReq) *base.Result {
	if req == nil || len(req.GetIds()) <= 0 {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	result := _self.GetDB().Delete(&entity.SysDictData{}, req.GetIds())
	if result.Error != nil {
		return base.ResultFailureErr(result.Error)
	}
	if result.RowsAffected <= 0 {
		_self.GetLogger().Error("删除失败,无删除记录")
		return base.ResultFailureMsg("删除失败")
	}
	return base.ResultSuccessMsg("删除成功")

}

func (_self *SysDictData) DictDataList(dictType *string) *base.Result {
	if dictType == nil || len(*dictType) <= 0 {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	var slice []vo.DictDataList
	err := _self.GetDB().Model(&entity.SysDictData{}).Order("sort").
		Find(&slice, "type_code = ? and status = '0'", dictType).Error
	if err != nil {
		return base.ResultFailureErr(err)
	}
	return base.ResultSuccess(slice)
}
