package dto

import "go-protector/server/internal/base"

type ApproveCmdPageReq struct {
	base.Pagination
	Cmd     string `json:"cmd,omitempty"`
	AssetId uint64 `json:"assetId,omitempty" binding:"required"`
}
