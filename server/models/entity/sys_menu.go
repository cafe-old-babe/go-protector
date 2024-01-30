package entity

import "go-protector/server/core/consts/table_name"

type SysMenu struct {
	ModelId
	MenuName   string `json:"menuName" gorm:"comment:菜单名称"`
	MenuType   int8   `json:"menuType" gorm:"comment:菜单类型,0-目录,1-菜单,2-按钮"`
	PID        uint64 `json:"PID" gorm:"comment:父级ID"`
	Permission string `json:"permission" gorm:"comment:权限标识"`
	Hidden     int8   `json:"hidden" gorm:"comment:是否隐藏,menu_type!=2时必填,1-隐藏,其他-显示"`
	RouterPath string `json:"routerPath" gorm:"comment:路由path,动态路由"`
	ModelDelete
	ModelControl
}

func (SysMenu) TableName() string {
	return table_name.SysMenu
}
