package dao

import (
	"errors"
	"go-protector/server/commons/local/table_name"
	"go-protector/server/commons/selfErr"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"gorm.io/gorm"
	"time"
)

var SysUser sysUser

type sysUser struct {
}

// FindUserByDTO 根据条件查询用户信息
func (_self *sysUser) FindUserByDTO(db *gorm.DB, dto *dto.FindUserDTO) (
	sysUser *entity.SysUser, err error) {

	if dto == nil || (len(dto.LoginName) <= 0 && dto.Id <= 0) {
		err = selfErr.ErrParamInvalid
		return
	}

	sysUser = new(entity.SysUser)

	tx := db.Scopes(func(db *gorm.DB) *gorm.DB {
		db.Where("login_name = ?", dto.LoginName)
		if dto.UserStatus == 0 {
			db.Where("user_status = ?", dto.UserStatus)
		}
		return db
	})
	if dto.IsUnscoped {
		tx = tx.Unscoped()
	}
	if err = tx.First(sysUser).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			sysUser = nil
			err = selfErr.ErrRecordNotFoundSysUser
			return
		}
	}
	return
}

// LockUser 锁定用户
// id
// lockType 锁定类型
// lockReason 锁定原因
func (_self *sysUser) LockUser(db *gorm.DB, id uint64, lockType int, lockReason string) (err error) {

	if id <= 0 || lockType <= 0 || len(lockReason) <= 0 {
		err = selfErr.ErrParamInvalid
		return
	}
	result := db.Table(table_name.SysUser).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"lock_type":   lockType,
			"lock_reason": lockReason,
		})
	if result.Error != nil {
		err = result.Error
		return
	}
	if result.RowsAffected <= 0 {
		err = selfErr.ErrUpdateFailure
	}

	return
}

// UpdateLastLoginIp 更新最后登录IP
func (_self *sysUser) UpdateLastLoginIp(db *gorm.DB, id uint64, lastLoginIp string) (err error) {
	if id <= 0 || len(lastLoginIp) <= 0 {
		err = selfErr.ErrParamInvalid
	}
	result := db.Table(table_name.SysUser).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_ip":   lastLoginIp,
			"last_login_time": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		err = selfErr.ErrUpdateFailure
	}
	return

}
