package entity

import (
	"go-protector/server/internal/approve"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/utils"
	"gorm.io/gorm"
)

type ApproveRecord struct {
	ModelId
	SessionId         string               `gorm:"column:session_id;size:64;comment:sessionId" json:"sessionId,omitempty"`
	ApplicantId       uint64               `gorm:"column:applicant_id;comment:申请人ID" json:"applicantId,omitempty"`
	ApplicantContent  string               `gorm:"column:applicant_content;comment:申请内容" json:"applicantContent,omitempty"`
	ApproveUserId     uint64               `gorm:"column:approve_user_id;comment:审批人Id" json:"approveUserId,omitempty"`
	ApproveContent    string               `gorm:"column:approve_content;comment:审批人回复" json:"approveContent,omitempty"`
	ApproveStatus     c_type.ApproveStatus `gorm:"column:approve_status;comment:审批记录状态" json:"approveStatus"`
	ApproveType       c_type.ApproveType   `gorm:"column:approve_type;comment:审批类型" json:"approveType,omitempty"`
	ApproveBindId     uint64               `gorm:"column:approve_bind_id;comment:审批记录关联ID" json:"approveBindId,omitempty"`
	Timeout           int                  `gorm:"column:timeout;comment:超时时间(秒)" json:"timeout"`
	WorkNum           string               `gorm:"column:work_num;comment:工单号"  json:"workNum" binding:"omitempty"`
	ApproveTypeName   string               `gorm:"-" json:"approveTypeName"`
	ApproveStatusText string               `gorm:"-" json:"approveStatusText"`
	ApplicantUser     SysUser              `gorm:"foreignKey:ApplicantId"  json:"applicantUser" binding:"omitempty"`
	ApproveUser       SysUser              `gorm:"foreignKey:ApproveUserId"  json:"approveUser" binding:"omitempty"`
	ModelControl
	ModelDelete
}

func (*ApproveRecord) TableName() string {
	return table_name.ApproveRecord
}

func (_self *ApproveRecord) BeforeCreate(db *gorm.DB) (err error) {
	if len(_self.WorkNum) <= 0 {
		_self.WorkNum, err = utils.GenerateDateNextSeq("approve-")
	}
	return
}

func (_self *ApproveRecord) AfterFind(db *gorm.DB) error {
	_self.ApproveStatusText = _self.ApproveStatus.String()
	info := approve.GetTypeInfo(_self.ApproveType)
	if info != nil {
		_self.ApproveTypeName = info.ApproveTypeText
	}
	_self.ApproveUser.Password = ""
	_self.ApplicantUser.Password = ""
	return nil
}

func (_self *ApproveRecord) AfterUpdate(db *gorm.DB) error {
	return _self.AfterFind(db)
}
