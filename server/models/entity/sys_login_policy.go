package entity

import (
	"go-protector/server/core/consts/table_name"
	"go-protector/server/core/custom/c_type"
)

type SysLoginPolicy struct {
	ModelId
	PolicyCode c_type.LoginPolicyCode `gorm:"size:32;comment:策略CODE"`
	PolicyName string                 `gorm:"size:32;comment:策略名称"`
	Enable     int                    `gorm:"size:1;comment:是否启用;1-启用;0-停用"`
	Json       string                 `gorm:"type:text;comment:配置参数"`
	Remark     string                 `gorm:"size:512;comment:策略说明"`
	ModelControl
}

func (SysLoginPolicy) TableName() string {
	return table_name.SysLoginPolicy
}
