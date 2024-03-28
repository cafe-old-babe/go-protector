package dto

import "go-protector/server/core/base"

type SysPostPageReq struct {
	base.Pagination
	Name    string   `json:"name"`
	DeptIds []uint64 `json:"deptIds"`
	Code    string   `json:"code"`
}

type SysPostSaveReq struct {
	ID      uint64   `json:"id"`
	DeptIds []uint64 `json:"deptIds" binding:"required"`
	Name    string   `json:"name" binding:"required"`
	Code    string   `json:"code" binding:"required"`
}
