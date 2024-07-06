package dto

import (
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_type"
)

type SsoOperationPageReq struct {
	base.Pagination
	AssetName    string `json:"assetName,omitempty"`
	AssetIp      string `json:"assetIp,omitempty"`
	AssetAcc     string `json:"assetAcc,omitempty"`
	Cmd          string `json:"cmd,omitempty"`
	SsoSessionId uint64 `json:"ssoSessionId,omitempty"`
}

type SsoOperationBySsoIdPage struct {
	ID           string      `json:"id,omitempty"`
	Cmd          string      `json:"cmd,omitempty"`
	TimeStamp    float64     `json:"timeStamp"`
	CmdStartAt   c_type.Time `json:"-"`
	SsoSessionId uint64      `json:"ssoSessionId,omitempty"`
}
