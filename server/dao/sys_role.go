package dao

import (
	"go-protector/server/core/consts/table_name"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_type"
	"go-protector/server/models/entity"
	"gorm.io/gorm"
)

var SysRole sysRole

type sysRole struct {
}

// GetRoleIdByRelationId 根据关联id查询角色id
func (_self sysRole) GetRoleIdByRelationId(db *gorm.DB, relationId uint64,
	relationType c_type.RoleRelationType) (roleIdSlice []uint64, err error) {
	if relationId <= 0 || len(relationType) <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	var relationSlice []entity.SysRoleRelation
	if err = db.Where("relation_id = ? and relation_type = ?", relationId, relationType).
		Find(&relationSlice).Error; err != nil {
		return
	}
	for _, relation := range relationSlice {
		roleIdSlice = append(roleIdSlice, relation.ID)
	}
	return
}

// GetPermissionSliceByIds 根据角色ID的权限 查询菜单列表
// ids 角色ID
// menuType 菜单类型
// admin 是否为管理员 非管理员:关联权限表 管理员:直接查询
func (_self sysRole) GetPermissionSliceByIds(db *gorm.DB, ids []uint64, menuType []int8, admin bool) (
	menuSlice []entity.SysMenu, err error) {
	tx := db.Debug().Table(table_name.SysMenu)

	if !admin {
		if len(ids) <= 0 {
			err = c_error.ErrParamInvalid
			return
		}
		// 子查询
		tx = tx.Where("id in (?) ",
			// 权限表
			db.Table(table_name.SysRoleRelation).Select("relation_id").
				Scopes(func(db *gorm.DB) *gorm.DB {

					if len(menuType) > 0 {
						db = db.Where("relation_type in ?", menuType)
					}
					return db
				}).
				Where("role_id in ? ", ids),
		)
	} else {
		if len(menuType) > 0 {
			tx = tx.Where("menu_type in ?", menuType)
		}
	}

	err = tx.Find(&menuSlice).Error
	return
}
