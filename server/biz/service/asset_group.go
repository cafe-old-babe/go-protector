package service

import (
	"errors"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_error"
	"gorm.io/gorm"
	"strings"
)

type AssetGroup struct {
	base.Service
}

func (_self *AssetGroup) Tree() *base.Result {
	var groupSlice []entity.AssetGroup
	if err := _self.GetDB().Find(&groupSlice).Error; err != nil {
		return base.ResultFailureErr(err)
	}
	node := dto.GenerateTree(groupSlice, 0, "ID", "PID", "Name", nil)
	return base.ResultSuccess(node)
}

func (_self *AssetGroup) SaveCheck(model *entity.AssetGroup) error {
	if model.ID == model.PID {
		return errors.New("父节点校验失败,请选择有效父节点")
	}
	var count int64
	err := _self.GetDB().Model(model).Scopes(func(db *gorm.DB) *gorm.DB {
		if model.ID > 0 {
			db = db.Where("id <> ?", model.ID)
		}

		return db.Where("name = ? ", model.Name)
	}).Where("p_id = ?", model.PID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("同级别资源组中不能同时存在:" + model.Name)
	}
	return nil
}

func (_self *AssetGroup) DeleteByIds(req *base.IdsReq) (result *base.Result) {
	if req == nil || len(req.GetIds()) <= 0 {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	// 校验资源组子集下是否有资产信息
	ids := req.GetIds()
	childrenIdSlice, err := _self.findChildrenIdSlice(ids, false)
	if err != nil {
		return base.ResultFailureErr(err)
	}
	if len(childrenIdSlice) > 0 {
		result = base.ResultFailureMsg("删除的资源组中包含资产信息,请移除资源组下的资产信息")
		return
	}
	// 删除所有
	if err = _self.GetDB().Delete(&entity.AssetGroup{}, ids).Error; err != nil {
		return base.ResultFailureErr(err)
	}
	return base.ResultSuccessMsg("删除成功")
}

// findChildrenIdSlice 查询子节点id
// contains 是否包含自己
func (_self *AssetGroup) findChildrenIdSlice(ids []uint64, contains bool) (idSlice []uint64, err error) {
	if ids == nil || len(ids) <= 0 {
		return
	}
	builder := strings.Builder{}
	builder.WriteString(" with recursive asset_group_rec as ( ")
	builder.WriteString(" select id,p_id ")
	builder.WriteString(" from asset_group where ")
	if contains {
		builder.WriteString(" id in (?) ")
	} else {
		builder.WriteString(" p_id in (?) ")
	}
	builder.WriteString(" union all ")
	builder.WriteString(" select ag.id,ag.p_id ")
	builder.WriteString(" from asset_group ag ")
	builder.WriteString(" join asset_group_rec agr on ag.p_id = agr.id ")
	builder.WriteString(" ) select * from asset_group_rec")
	subQuery := _self.GetDB().Raw(builder.String(), ids)
	err = _self.GetDB().Table("(?) as t", subQuery).
		Pluck("id", &idSlice).Error

	return
}
