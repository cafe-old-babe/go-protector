package entity

import (
	"database/sql"
	"go-protector/server/internal/consts/table_name"
)

type SysMenu struct {
	ModelId
	Name       string        `json:"name" gorm:"size:64;comment:菜单/按钮名称"`
	MenuType   int8          `json:"menuType" gorm:"comment:菜单类型,0-目录,1-菜单,2-按钮"`
	PID        uint64        `json:"PID" gorm:"comment:父级ID"`
	Permission string        `json:"permission" gorm:"size:64;comment:权限标识"`
	Hidden     sql.NullInt16 `json:"hidden" gorm:"comment:是否隐藏,menu_type!=2时必填,1-隐藏,其他-显示"`
	Path       string        `json:"path" gorm:"size:64;comment:path"`
	Component  string        `json:"component" gorm:"size:64;comment:组件名称"`
	ModelDelete
	ModelControl
}

func (SysMenu) TableName() string {
	return table_name.SysMenu
}
