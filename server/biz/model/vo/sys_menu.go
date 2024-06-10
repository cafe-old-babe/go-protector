package vo

type SysMenuVO struct {
	ID           uint64      `json:"id"`
	PID          uint64      `json:"pid"`
	Name         string      `json:"name"`
	Path         string      `json:"path"`
	MenuType     int8        `json:"menuType"`
	MenuTypeName string      `json:"menuTypeName"`
	Permission   string      `json:"permission"`
	Hidden       int16       `json:"hidden"`
	Component    string      `json:"component"`
	Children     []SysMenuVO `json:"children"`
}
