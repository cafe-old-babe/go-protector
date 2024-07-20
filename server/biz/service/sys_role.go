package service

import (
	"errors"
	"go-protector/server/biz/model/dao"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
	"gorm.io/gorm"
)

type SysRole struct {
	base.Service
}

// GetMenuByRoleIds 根据角色获取菜单权限
func (_self *SysRole) GetMenuByRoleIds(roleIds []uint64, isAdminParam ...bool) (
	menuSlice, buttonSlice []entity.SysMenu, err error) {
	var isAdmin bool
	if len(isAdminParam) > 0 && isAdminParam[0] {
		isAdmin = isAdminParam[0]
	}

	menuSlice, err = dao.SysRole.GetPermissionSliceByIds(_self.GetDB(), roleIds, []int8{0, 1}, isAdmin)
	if err != nil {
		return
	}
	buttonSlice, err = dao.SysRole.GetPermissionSliceByIds(_self.GetDB(), roleIds, []int8{2}, isAdmin)

	return
}

func (_self *SysRole) Page(req *dto.SysRolePageReq) (res *base.Result) {
	if req == nil {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	var slice []entity.SysRole
	var count int64
	if err := _self.GetDB().Scopes(
		condition.Paginate(req.GetPagination()),
		condition.Like("role_name", req.RoleName),
	).Find(&slice).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		res = base.ResultFailureErr(err)
	} else {
		res = base.ResultPage(slice, req.GetPagination(), count)
	}
	return
}

func (_self *SysRole) SaveCheck(entity *entity.SysRole) (err error) {
	var count int64
	err = _self.GetDB().Model(entity).Scopes(func(db *gorm.DB) *gorm.DB {
		if entity.ID > 0 {
			db = db.Where("id <> ?", entity.ID)
		}
		return db.Where("role_name = ? ", entity.RoleName)
	}).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		err = errors.New("角色名称不能重复:" + entity.RoleName)
	}
	return nil
}

// Delete 删除
func (_self *SysRole) Delete(req *base.IdsReq) *base.Result {
	if err := dao.SysRole.DeleteByRoleIds(_self.GetDB(), req.GetIds()); err != nil {
		return base.ResultFailureErr(err)
	}
	return base.ResultSuccessMsg("删除成功")
}

// GetPermission 获取权限
func (_self *SysRole) GetPermission(roleId uint64) *base.Result {
	if roleId <= 0 {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}

	var menuIdSlice []uint64
	if err := _self.GetDB().Table(table_name.SysRoleRelation).
		Where("role_id = ? and relation_type = ?", roleId, consts.Menu).
		Pluck("relation_id", &menuIdSlice).Error; err != nil {
		return base.ResultFailureErr(err)
	}

	return base.ResultSuccess(menuIdSlice)
}

// SavePermission 保存权限
func (_self *SysRole) SavePermission(roleId uint64, ids []uint64) (res *base.Result) {
	err := dao.SysRole.SavePermission(_self.GetDB(), roleId, ids)
	if err != nil {
		res = base.ResultFailureErr(err)
	} else {
		res = base.ResultSuccessMsg("保存成功")
	}
	return
}

func (_self *SysRole) SetStatus(roleId uint64, status int8) *base.Result {
	// check
	if status != 0 && status != 1 {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	var sysRole entity.SysRole
	//没有找到记录时，它会返回 ErrRecordNotFound 错误
	if err := _self.GetDB().Take(&sysRole, roleId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return base.ResultSuccessMsg("无效更新")
		}
		return base.ResultFailureErr(err)
	}
	if sysRole.IsInner == 1 {
		return base.ResultFailureMsg("无法修改内置角色的状态")
	}

	if sysRole.Status == status {
		return base.ResultSuccessMsg("未改变状态")
	}
	update := _self.GetDB().Model(&sysRole).
		Where("id = ? and status <> ?", roleId, status).
		Update("status", status)
	if update.Error != nil {
		return base.ResultFailureErr(update.Error)
	}
	if update.RowsAffected <= 0 {
		return base.ResultFailureMsg("更新失败")
	}
	return base.ResultSuccessMsg("更新成功")
}

func (_self *SysRole) List() (res *base.Result) {
	var slice []entity.SysRole
	if err := _self.GetDB().Find(&slice).Error; err != nil {
		res = base.ResultFailureErr(err)
	} else {
		res = base.ResultSuccess(slice)
	}
	return
}
