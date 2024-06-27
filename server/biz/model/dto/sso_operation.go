package dto

import "go-protector/server/internal/base"

type SsoOperationPageReq struct {
	base.Pagination
	AssetName string `json:"assetName,omitempty"`
	AssetIp   string `json:"assetIp,omitempty"`
	AssetAcc  string `json:"assetAcc,omitempty"`
	Cmd       string `json:"cmd,omitempty"`
}
