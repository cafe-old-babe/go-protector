package base

import "gorm.io/gorm/schema"

type IdsReq struct {
	Ids      []uint64      `json:"ids" binding:"required"`
	Value    schema.Tabler `json:"-"`
	Unscoped bool          `json:"-"`
}

// GetIds 获取ids
func (_self IdsReq) GetIds() []uint64 {
	return _self.Ids
}
