package dto

type RoleInfo struct {
	Name string `json:"name"`
	Role Role   `json:"role"`
}

type Permission struct {
	PermissionId    string         `json:"permissionId"`
	PermissionName  string         `json:"permissionName"`
	Actions         string         `json:"actions"`
	ActionEntitySet []ActionEntity `json:"actionEntitySet"`
}

type ActionEntity struct {
	Action       string `json:"action"`
	Describe     string `json:"describe"`
	DefaultCheck bool   `json:"defaultCheck"`
}

type Role struct {
	Permissions []Permission `json:"permissions"`
}

type MenuInfo struct {
	Id        uint64      `json:"id"`
	ParentId  uint64      `json:"parentId"`
	Name      string      `json:"name"`      // 路由名称 必填不能重复
	Path      string      `json:"path"`      // 路由url路径 非必填
	Component string      `json:"component"` // 组件名称 必填
	Redirect  interface{} `json:"redirect"`
	Meta      MetaInfo    `json:"meta"`
}

type MetaInfo struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Show  bool   `json:"show"`
}

type SysRolePageReq struct {
	Pagination `json:"-"`
	RoleName   string `json:"roleName"`
}
