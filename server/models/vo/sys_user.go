package vo

type UserPage struct {
	ID        uint64 `json:"id"`
	LoginName string `json:"loginName"`
	Username  string `json:"username"`
	Sex       string `json:"sex"`
	Status    int    `json:"status"`
	DeptId    int    `json:"deptId"`
	DeptName  string `json:"deptName"`
}
