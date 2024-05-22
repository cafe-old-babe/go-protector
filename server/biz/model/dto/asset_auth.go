package dto

import "go-protector/server/internal/base"

type AssetAuthPageReq struct {
	base.Pagination
	AssetIp  string `json:"assetIp,omitempty"`
	AssetAcc string `json:"assetAcc,omitempty"`
	UserAcc  string `json:"userAcc,omitempty"`
}
