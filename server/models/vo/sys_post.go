package vo

import "gorm.io/gorm"

type PostPage struct {
	ID        uint64         `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	DeptIds   string         `json:"deptIds"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
