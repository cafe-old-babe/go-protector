package entity

import (
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_type"
)

// 关联 user_id, menu_id, dept_id

type SysRoleRelation struct {
	ModelId
	RoleId       uint64              `gorm:"comment:角色ID"`
	RelationId   uint64              `gorm:"comment:关联ID"`
	RelationType c_type.RelationType `gorm:"comment:关联类型,user-用户,menu-菜单"`
}

func (*SysRoleRelation) TableName() string {
	return table_name.SysRoleRelation
}
