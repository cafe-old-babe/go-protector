package service

import (
	"go-protector/server/core/base"
	"go-protector/server/dao"
	"go-protector/server/models/entity"
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
