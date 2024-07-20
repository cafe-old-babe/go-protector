package entity

import (
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_type"
)

type ApproveRecord struct {
	ModelId
	SessionId        string               `gorm:"column:session_id;size:64;comment:sessionId" json:"sessionId,omitempty"`
	ApplicantId      uint64               `gorm:"column:applicant_id;comment:申请人ID" json:"applicantId,omitempty"`
	ApplicantContent string               `gorm:"column:applicant_content;comment:申请内容" json:"applicantContent,omitempty"`
	ApproveUserId    uint64               `gorm:"column:approve_user_id;comment:审批人Id" json:"approveUserId,omitempty"`
	ApproveContent   string               `gorm:"column:approve_content;comment:审批人回复" json:"approveContent,omitempty"`
	ApproveStatus    c_type.ApproveStatus `gorm:"column:approve_status;comment:审批记录状态" json:"approveStatus,omitempty"`
	ApproveType      c_type.ApproveType   `gorm:"column:approve_type;comment:审批类型" json:"approveType,omitempty"`
	ApproveBindId    uint64               `gorm:"column:approve_bind_id;comment:审批记录关联ID" json:"approveBindId,omitempty"`
	Timeout          int                  `gorm:"column:timeout;comment:超时时间(秒)" json:"timeout"`
	ApproveTypeName  string               `gorm:"-" json:"approveTypeName"`
	ApplicantUser    SysUser              `gorm:"foreignKey:ApplicantId"  json:"applicantUser" binding:"omitempty"`
	ApproveUser      SysUser              `gorm:"foreignKey:ApproveUserId"  json:"approveUser" binding:"omitempty"`
	ModelControl
	ModelDelete
}

func (_self ApproveRecord) TableName() string {
	return table_name.ApproveRecord

}
