package entity

import "go-protector/server/internal/consts/table_name"

type SysPost struct {
	ModelId
	Name string `json:"name" gorm:"size:64;comment:岗位名称" binding:"required"`
	Code string `json:"code" gorm:"size:64;comment:岗位代码" binding:"required"`
	ModelDelete
	ModelControl
}

func (SysPost) TableName() string {
	return table_name.SysPost
}
