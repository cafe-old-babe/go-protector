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
	AssetAcc  string `json:"assetAcc"`
}

type OperationForMonitorReq struct {
	SsoSessionId uint64 `json:"ssoSessionId,omitempty" binding:"required"`
	Type         string `json:"type,omitempty" binding:"required"` // 1: 告警;2: 阻断
	Message      string `json:"message,omitempty" binding:"required_if=Type 1"`
}
