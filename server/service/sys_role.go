package service

import (
	"errors"
	"go-protector/server/core/base"
	"go-protector/server/core/consts"
	"go-protector/server/core/consts/table_name"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/scope"
	"go-protector/server/dao"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
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

	menuSlice, err = dao.SysRole.GetPermissionSliceByIds(_self.DB, roleIds, []int8{0, 1}, isAdmin)
	if err != nil {
		return
	}
	buttonSlice, err = dao.SysRole.GetPermissionSliceByIds(_self.DB, roleIds, []int8{2}, isAdmin)

	return
}

func (_self *SysRole) Page(req *dto.SysRolePageReq) (res *dto.Result) {
	if req == nil {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
	}
	var slice []entity.SysRole
	var count int64
	if err := _self.DB.Scopes(
		scope.Paginate(req.GetPagination()),
		scope.Like("role_name", req.RoleName),
	).Find(&slice).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		res = dto.ResultFailureErr(err)
	} else {
		res = dto.ResultPage(slice, req.GetPagination(), count)
	}
	return
}

func (_self *SysRole) SaveCheck(entity *entity.SysRole) (err error) {
	var count int64
	err = _self.DB.Model(entity).Scopes(func(db *gorm.DB) *gorm.DB {
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
func (_self *SysRole) Delete(req *dto.IdsReq) *dto.Result {
	if err := dao.SysRole.DeleteByRoleIds(_self.DB, req.GetIds()); err != nil {
		return dto.ResultFailureErr(err)
	}
	return dto.ResultSuccessMsg("删除成功")
}

// GetPermission 获取权限
func (_self *SysRole) GetPermission(roleId uint64) *dto.Result {
	if roleId <= 0 {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
	}

	var menuIdSlice []uint64
	if err := _self.DB.Table(table_name.SysRoleRelation).
		Where("role_id = ? and relation_type = ?", roleId, consts.Menu).
		Pluck("relation_id", &menuIdSlice).Error; err != nil {
		return dto.ResultFailureErr(err)
	}

	return dto.ResultSuccess(menuIdSlice)
}

// SavePermission 保存权限
func (_self *SysRole) SavePermission(roleId uint64, ids []uint64) (res *dto.Result) {
	err := dao.SysRole.SavePermission(_self.DB, roleId, ids)
	if err != nil {
		res = dto.ResultFailureErr(err)
	} else {
		res = dto.ResultSuccessMsg("保存成功")
	}
	return
}

func (_self *SysRole) SetStatus(roleId uint64, status int8) *dto.Result {
	// check
	if status != 0 && status != 1 {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
	}
	var sysRole entity.SysRole
	//没有找到记录时，它会返回 ErrRecordNotFound 错误
	if err := _self.DB.Take(&sysRole, roleId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ResultSuccessMsg("无效更新")
		}
		return dto.ResultFailureErr(err)
	}
	if sysRole.IsInner == 1 {
		return dto.ResultFailureMsg("无法修改内置角色的状态")
	}

	if sysRole.Status == status {
		return dto.ResultSuccessMsg("未改变状态")
	}
	update := _self.DB.Model(&sysRole).
		Where("id = ? and status <> ?", roleId, status).
		Update("status", status)
	if update.Error != nil {
		return dto.ResultFailureErr(update.Error)
	}
	if update.RowsAffected <= 0 {
		return dto.ResultFailureMsg("更新失败")
	}
	return dto.ResultSuccessMsg("更新成功")
}
