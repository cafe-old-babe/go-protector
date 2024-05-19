package dto

import (
	"go-protector/server/internal/base"
)

type AssetGatewayPageReq struct {
	base.Pagination
	AgName string `json:"agName"`
	AgIp   string `json:"agIp"`
}
