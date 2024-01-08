package dao

import (
	"database/sql"
	"errors"
	"go-protector/server/core/current"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/local/table_name"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"gorm.io/gorm"
	"time"
)

var SysUser sysUser

type sysUser struct {
}

// FindUserByDTO 根据条件查询用户信息
func (_self *sysUser) FindUserByDTO(db *gorm.DB, dto *dto.FindUser) (
	sysUser *entity.SysUser, err error) {

	if dto == nil || (len(dto.LoginName) <= 0 && dto.ID <= 0) {
		err = c_error.ErrParamInvalid
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
			err = c_error.ErrRecordNotFoundSysUser
			return
		}
	}
	return
}

// LockUser 锁定用户
// entity.SysUser 锁定用户信息
func (_self *sysUser) LockUser(db *gorm.DB, entity *entity.SysUser) (err error) {

	if nil == entity || entity.ID <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	selectSlice := []string{"user_status", "lock_reason", "lock_time"}

	if entity.UpdatedBy <= 0 {
		if userId := current.GetUserId(db.Statement.Context); userId > 0 {
			entity.UpdatedBy = userId
			selectSlice = append(selectSlice, "updated_by")
		}
	}

	entity.LockTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	// 防止ABA的问题
	result := db.Model(entity).Where("user_status = ?", 0).Select(selectSlice).
		Updates(entity)

	if result.Error != nil {
		err = result.Error
		return
	}
	if result.RowsAffected <= 0 {
		err = c_error.ErrUpdateFailure
	}

	return
}

// UnlockUser 解锁用户
func (_self *sysUser) UnlockUser(db *gorm.DB, dto *dto.SetStatus) error {
	if nil == dto || dto.ID <= 0 {
		return c_error.ErrParamInvalid
	}

	updateMap := map[string]interface{}{
		"lock_time":   nil,
		"lock_reason": nil,
		"user_status": 0,
	}
	if len(dto.ExpirationAt) > 0 {
		parse, err := time.Parse(dto.ExpirationAt, time.DateTime)
		if err != nil {
			return err
		}
		updateMap["expiration_at"] = parse
	}
	res := db.Table(table_name.SysUser).
		Where("id = ? and user_status != 0", dto.ID). // 防止ABA的问题
		Updates(updateMap)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected <= 0 {
		return c_error.ErrUpdateFailure
	}
	return nil
}

// UpdateLastLoginIp 更新最后登录IP
func (_self *sysUser) UpdateLastLoginIp(db *gorm.DB, id uint64, lastLoginIp string) (err error) {
	if id <= 0 || len(lastLoginIp) <= 0 {
		err = c_error.ErrParamInvalid
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
		err = c_error.ErrUpdateFailure
	}
	return

}
