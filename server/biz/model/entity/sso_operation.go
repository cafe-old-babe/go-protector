package entity

import (
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_type"
	"gorm.io/gorm"
)

type SsoOperation struct {
	ModelId
	SsoSessionId uint64      `gorm:"comment:会话ID" json:"ssoSessionId"`
	Sort         int         `gorm:"comment:排序" json:"sort"`
	PS1          string      `gorm:"size:64;comment:PS1" json:"PS1"`
	Cmd          string      `gorm:"size:4096;comment:执行命令" json:"cmd"`
	ExecStatus   string      `gorm:"size:1;comment:执行状态,0:执行成功,1:待审批,2:审批未通过" json:"execStatus"`
	CmdStartAt   c_type.Time `gorm:"comment:命令开始时间" json:"cmdStartAt"`
	CmdExecAt    c_type.Time `gorm:"comment:命令执行时间" json:"cmdExecAt"`
	SsoSession   SsoSession  `gorm:"foreignKey:SsoSessionId"  json:"ssoSession,omitempty" binding:"-"`
	ModelControl
	ModelDelete
}

func (_self *SsoOperation) TableName() string {
	return table_name.SsoOperation
}

func (_self *SsoOperation) AfterFind(db *gorm.DB) (err error) {

	_self.SsoSession.Completion()
	return
}
