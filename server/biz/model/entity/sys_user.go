package entity

import (
	"database/sql"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_type"
	"gorm.io/gorm"
)

type SysUser struct {
	ModelId
	LoginName     string         `json:"loginName" gorm:"size:32;comment:登录名"`
	Password      string         `json:"password" gorm:"size:256;comment:密码"`
	Username      string         `json:"username" gorm:"size:8;comment:用户名"`
	Email         string         `json:"email"  gorm:"size:128;comment:邮箱"`
	DeptId        uint64         `json:"deptId" gorm:"comment:部门ID"`
	Sex           string         `json:"sex" gorm:"size:32;comment:性别,0:女,1:男"`
	LastLoginTime sql.NullTime   `json:"lastLoginTime" gorm:"comment:最后登录时间"`
	LastLoginIp   string         `json:"lastLoginIp" gorm:"size:32;comment:最后登录IP"`
	LockTime      sql.NullTime   `json:"lockTime" gorm:"comment:锁定时间"`
	LockReason    sql.NullString `json:"lockReason" gorm:"size:256;comment:锁定原因"`
	UserStatus    int            `json:"userStatus"  gorm:"size:1;comment:用户状态,0:正常,非零锁定"` // 0: 正常 非零锁定,
	ExpirationAt  c_type.Time    `json:"expirationAt" gorm:"comment:有效时间"`
	SysRoleIds    []uint64       `json:"sysRoleIds" gorm:"-"`                            //
	SysRoles      []SysRole      `json:"sysRoles" gorm:"-"`                              //
	SysDept       SysDept        `json:"sysDept" gorm:"foreignKey:DeptId;references:ID"` // https://gorm.io/zh_CN/docs/belongs_to.html#Belongs-to-%E7%9A%84-CRUD
	SessionId     string         `json:"sessionId,omitempty" gorm:"-"`
	ModelDelete
	ModelControl
}

func (_self SysUser) BeforeUpdate(db *gorm.DB) error {
	var auth AssetAuth
	return auth.UpdateRedundancy(db, _self)
}

func (_self SysUser) TableName() string {
	return table_name.SysUser
}
