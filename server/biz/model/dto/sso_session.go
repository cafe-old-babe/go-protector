package dto

import "go-protector/server/internal/base"

type ConnectBySessionReq struct {
	Id uint64
	H  int `form:"h"`
	W  int `form:"w"`
}

type SsoSessionPageReq struct {
	base.Pagination
	AssetName string `json:"assetName,omitempty"`
	AssetIp   string `json:"assetIp,omitempty"`
	UserAcc   string `json:"userAcc"`
}
