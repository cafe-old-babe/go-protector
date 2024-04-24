package dao

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
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
		return db
	})
	if dto.IsUnscoped {
		tx = tx.Unscoped()
	}
	if err = tx.Preload("SysDept").First(sysUser).Error; err != nil {

		return
		//if errors.Is(gorm.ErrRecordNotFound, err) {
		//	sysUser = nil
		//	err = c_error.ErrRecordNotFoundSysUser
		//	return
		//}
	}

	var roleIdSlice []uint64
	if roleIdSlice, err = SysRole.GetRoleIdByRelationId(db, sysUser.ID, consts.User); err != nil {
		return
	}
	if len(roleIdSlice) <= 0 {
		err = c_error.ErrNotFoundRoleOfSysUser
	}
	var sysRoles []entity.SysRole
	if err = db.Find(&sysRoles, roleIdSlice).Error; err == nil {
		sysUser.SysRoles = sysRoles
		sysUser.SysRoleIds = roleIdSlice
	}
	if len(sysRoles) <= 0 {
		err = c_error.ErrNotFoundRoleOfSysUser
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
	// 使用 struct 更新时, GORM 将只更新非零值字段。 你可能想用 map 来更新属性，或者使用 Select 声明字段来更新
	updateMap := map[string]interface{}{
		"lock_time":     nil,
		"lock_reason":   nil,
		"user_status":   0,
		"expiration_at": nil,
	}
	//if len(dto.ExpirationAt) > 0 {
	//	parse, err := time.Parse(time.DateTime, dto.ExpirationAt)
	//	if err != nil {
	//		return err
	//	}
	//	updateMap["expiration_at"] = parse
	//}
	updateMap["expiration_at"] = dto.ExpirationAt
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

func (_self *sysUser) Save(db *gorm.DB, req *dto.UserSaveReq) (err error) {
	// 校验
	if req == nil {
		return c_error.ErrParamInvalid
	}
	if err = binding.Validator.ValidateStruct(req); err != nil {
		return
	}
	model := &entity.SysUser{
		ModelId: entity.ModelId{
			ID: req.ID,
		},
		LoginName:    req.LoginName,
		Password:     req.Password,
		ExpirationAt: req.ExpirationAt,
		Email:        req.Email,
		Username:     req.Username,
		DeptId:       req.DeptId,
		Sex:          req.Sex,
	}
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if model.ID <= 0 {
		var count int64
		if err = db.Model(model).Where("login_name = ?",
			req.LoginName).Count(&count).Error; err != nil {
			return
		}
		if count > 0 {
			return errors.New("用户帐号重复,请核对")
		}
		// 新增
		if err = tx.Create(model).Error; err != nil {
			return
		}
	} else {
		// 更新
		// Updates 方法支持 struct 和 map[string]interface{} 参数。当使用 struct 更新时，默认情况下GORM 只会更新非零值的字段
		if err = tx.Omit(
			"password",
			"login_name",
		).Updates(model).Error; err != nil {
			return
		}
	}
	// 绑定角色
	if err = SysRole.UserIdBindRoleIds(tx, model.ID, req.RoleIds); err != nil {
		return
	}
	// 绑定岗位
	err = SysPost.UserBindPostIds(tx, model.ID, req.PostIds)

	return
}
