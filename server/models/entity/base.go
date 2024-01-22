package entity

import (
	"database/sql"
	"gorm.io/gorm"
)

// ModelId https://gorm.io/zh_CN/docs/models.html
type ModelId struct {
	ID uint64 `json:"id" gorm:"primaryKey"`
}

type ModelControl struct {
	CreatedAt    sql.NullTime `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
	CreatedBy    uint64       `json:"createdBy" gorm:"comment:创建人"`
	UpdatedAt    sql.NullTime `json:"updatedAt" gorm:"autoUpdateTime;comment:修改时间"`
	UpdatedBy    uint64       `json:"updatedBy" gorm:"comment:修改人"`
	CreateByName string       `json:"createByName" gorm:"-"`
	UpdateByName string       `json:"updateByName" gorm:"-"`
}

type ModelDelete struct {
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除"`
}
