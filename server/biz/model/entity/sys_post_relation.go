package entity

import (
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_type"
)

type SysPostRelation struct {
	ModelId
	PostId       uint64              `gorm:"index;comment:岗位ID"`
	RelationId   uint64              `gorm:"comment:关联ID"`
	RelationType c_type.RelationType `gorm:"index;comment:关联类型,user-用户,dept-部门"`
}

func (*SysPostRelation) TableName() string {
	return table_name.SysPostRelation
}
