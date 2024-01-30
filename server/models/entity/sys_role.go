package entity

import "go-protector/server/core/consts/table_name"

type SysRole struct {
	ModelId
	RoleName string `json:"roleName" gorm:"size:32;comment:角色名称"`
	RoleType int8   `json:"roleType" gorm:"comment:角色类型,0-管理角色,1-普通角色"`
	Status   int8   `json:"status" gorm:"comment:角色状态,0-正常,1-停用"`
	Sort     int32  `json:"sort" gorm:"comment:排序"`
	Remark   string `json:"remark" gorm:"size:1024;comment:备注"`
	Inner    int8   `json:"inner" gorm:"comment:是否内置角色,1-是"`
	ModelControl
	ModelDelete
}

func (SysRole) TableName() string {
	return table_name.SysRole
}
