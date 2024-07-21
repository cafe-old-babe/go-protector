package dto

import (
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_type"
)

type ApproveRecordPageReq struct {
	base.Pagination
	WorkNum           string `json:"workNum,omitempty"`
	ApproveUsername   string `json:"approveUsername,omitempty"`
	ApplicantUsername string `json:"applicantUsername,omitempty"`
}

type ApproveRecordInsertDTO struct {
	Id               uint64             `json:"id,omitempty" binding:"required"`
	ApplicantId      uint64             `json:"applicantId,omitempty" binding:"required"`
	ApproveUserId    uint64             `json:"approveUserId,omitempty" binding:"required"`
	SessionId        string             `json:"sessionId,omitempty" binding:"required"`
	ApplicantContent string             `json:"applicantContent,omitempty" binding:"required"`
	Timeout          int                `json:"timeout,omitempty"`
	ApproveType      c_type.ApproveType `json:"approveType,omitempty" binding:"required"`
	ApproveBindId    uint64             `json:"approveBindId,omitempty" binding:"required"`
}

type DoApproveDTO struct {
	Id             uint64               `json:"id,omitempty" binding:"required"`
	ApproveStatus  c_type.ApproveStatus `json:"approveStatus,omitempty" binding:"required"`
	ApproveContent string               `json:"approveContent,omitempty"`
	ApproveUserId  uint64               `json:"approveUserId"` // 优先获取currentUser
}

type ApproveRecordDTO struct {
}
