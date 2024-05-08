package dto

import "go-protector/server/internal/base"

type AssetInfoPageReq struct {
	base.Pagination
	GroupIds  []uint64 `json:"groupIds,omitempty"`
	AssetName string   `json:"assetName,omitempty"`
	IP        string   `json:"IP,omitempty"`
}

type AssetInfoSaveReq struct {
	ID             uint64 `json:"id"`
	AssetName      string `json:"assetName" binding:"required" `
	GroupId        uint64 `json:"groupId" binding:"required"`
	IP             string `json:"ip" binding:"required,ip"`
	Port           uint   `json:"port" binding:"required,min=1,max=65535"`
	AssetGatewayId uint64 `json:"gatewayId"`
	ManagerUserId  uint64 `json:"managerUserId" binding:"required"`
	Account        string `json:"account" binding:"required"`
	Password       string `json:"password" binding:"required_without=ID"`
	//PrivilegeAccount *AssetAccountAccessReq `json:"privilegeAccount" binding:"required"`
}
