package dto

import "go-protector/server/core/base"

type AssetNetworkPageReq struct {
	base.Pagination
	AnName string `json:"anName"`
	AnIp   string `json:"anIp"`
}
