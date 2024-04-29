package dto

import "go-protector/server/internal/base"

// AssetAccountAccessReq 接入从帐号
type AssetAccountAccessReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	// AccountType 从帐号类型,0-特权帐号,1-管理帐号(管理帐号可执行sudo),2-普通帐号(普通帐号不可执行sudo)
	AccountType string `json:"accountType" binding:"required,oneof=0 1 2"`
	//AssetId     uint64 ` json:"assetId" binding:"required"`
}

type AssetAccountPageReq struct {
	base.Pagination
	Account   string `json:"account,omitempty"`
	AssetName string `json:"assetName,omitempty"`
	IP        string `json:"ip,omitempty"`
}

type AccountAnalysisExtendDTO struct {
	ID      uint64
	AssetId uint64
	In      string
	Out     []byte
	Err     error
}
