package dto

type SysMenuSaveReq struct {
	ID         uint64 `json:"id"`
	PID        uint64 `json:"pid" binding:"required"`
	Name       string `json:"name" binding:"required,max=64"`
	Path       string `json:"path" binding:"max=64"`
	MenuType   int8   `json:"menuType" binding:"required,oneof=0 1 2"`
	Permission string `json:"permission" binding:"required,max=64"`
	Hidden     int8   `json:"hidden"`
	Component  string `json:"component"`
}
