package dao

import (
	"errors"
	"go-protector/server/core/consts"
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
	relationType c_type.RelationType) (roleIdSlice []uint64, err error) {
	if relationId <= 0 || len(relationType) <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	subQuery := db.Table(table_name.SysRole).Where("status <> 1").Select("id")
	var relationSlice []entity.SysRoleRelation
	if err = db.Where("relation_id = ? and relation_type = ?", relationId, relationType).
		Where("role_id in (?)", subQuery).
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
	subQuery := db.Table(table_name.SysRole).Where("id in ? and status <> 1", ids).Select("id")

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
				Where("role_id in (?) ", subQuery),
		)
	} else {
		if len(menuType) > 0 {
			tx = tx.Where("menu_type in ?", menuType)
		}
	}

	err = tx.Find(&menuSlice).Error
	return
}

// DeleteByRoleIds 同步删除
func (_self sysRole) DeleteByRoleIds(db *gorm.DB, ids []uint64) error {
	if len(ids) <= 0 {
		return c_error.ErrParamInvalid
	}
	return db.Transaction(func(tx *gorm.DB) error {
		roleRelationPtr := &entity.SysRoleRelation{}
		// 校验是否绑定了人员
		var count int64
		if err := db.Model(roleRelationPtr).
			Where("role_id = ? and relation_type = ?", ids, consts.User).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("删除的角色中绑定了用户,无法删除")
		}
		// 删除角色表
		if err := db.Unscoped().
			Where("id in ?", ids).
			Delete(&entity.SysRole{}).Error; err != nil {
			return err
		}
		// 删除关联关系
		return db.Where("role_id in ?", ids).Delete(roleRelationPtr).Error

	})
}

func (_self sysRole) SavePermission(db *gorm.DB, roleId uint64, ids []uint64) error {
	if roleId <= 0 || len(ids) <= 0 {
		return c_error.ErrParamInvalid
	}
	return db.Transaction(func(tx *gorm.DB) error {

		model := &entity.SysRoleRelation{}
		if err := tx.Delete(model,
			"role_id = ? and relation_type = ?", roleId, consts.Menu).Error; err != nil {
			return err
		}
		// 当使用map来创建时，钩子方法不会执行，关联不会被保存且不会回写主键。
		var relationSlice []map[string]interface{}
		for _, id := range ids {
			relationSlice = append(relationSlice, map[string]interface{}{
				"role_id":       roleId,
				"relation_type": consts.Menu,
				"relation_id":   id,
			})
		}
		return tx.Model(model).CreateInBatches(relationSlice, 100).Error
	})
}
