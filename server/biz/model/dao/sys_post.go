package dao

import (
	"errors"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/utils"
	"gorm.io/gorm"
	"strings"
)

var SysPost sysPost

type sysPost struct {
}

func (_self sysPost) DeleteByPostId(db *gorm.DB, ids []uint64) (err error) {
	if len(ids) <= 0 {
		return c_error.ErrParamInvalid
	}
	var count int64
	err = db.Table(table_name.SysPostRelation).
		Where("post_id in ? and relation_type = ?", ids, consts.User).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		err = errors.New("要删除的数据中,包含正在使用的岗位,无法删除")
	}

	return db.Transaction(func(tx *gorm.DB) (err error) {
		// 删除 sysPost
		if err = tx.Delete(&entity.SysPost{}, ids).Error; err != nil {
			return
		}
		// todo 删除授权关系
		return tx.Delete(&entity.SysPostRelation{}, "post_id in ? ", ids).Error

	})
}

func (_self sysPost) Save(db *gorm.DB, req *dto.SysPostSaveReq) error {
	var count int64
	if err := db.Table(table_name.SysPost).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if req.ID > 0 {
				db = db.Where("id <> ? ", req.ID)
			}
			return db
		}).Where("name = ? or LOWER(code) = ?",
		req.Name, strings.ToLower(req.Code)).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("岗位名称和岗位代码不能重复")
	}
	// 更新
	if req.ID > 0 {
		// 查询当前的部门关联关系
		var deptIds []uint64
		if err := db.Table(table_name.SysUser).Where(
			"id in (?)", db.Table(table_name.SysPostRelation).Select("relation_id").
				Where("relation_type = ? and post_id = ?", consts.User, req.ID)).
			Pluck("dept_id", &deptIds).Error; err != nil {
			return err
		}
		// deptIds 差集
		sub := utils.SliceSub(deptIds, req.DeptIds)
		if len(sub) > 0 {
			return errors.New("该岗位已被用户使用,无法删除部门信息,请核对")
		}
	}

	return db.Transaction(func(tx *gorm.DB) (err error) {
		var model entity.SysPost
		model = entity.SysPost{
			ModelId: entity.ModelId{ID: req.ID},
			Name:    req.Name,
			Code:    req.Code,
		}
		// 新增或更新
		if req.ID <= 0 {
			// 新增
			err = db.Create(&model).Error
		} else {
			err = db.Updates(&model).Error
		}
		if err != nil {
			return err
		}
		// todo 删除sub 授权关系

		// 删除 关联关系
		if err = db.Delete(&entity.SysPostRelation{},
			"post_id = ? and relation_type = ?", req.ID, consts.Dept).
			Error; err != nil {
			return err
		}
		// 新增 关联关系
		var entitySlice []entity.SysPostRelation

		for _, deptId := range req.DeptIds {
			entitySlice = append(entitySlice, entity.SysPostRelation{
				PostId:       model.ID,
				RelationId:   deptId,
				RelationType: consts.Dept,
			})
		}
		if err = db.Create(&entitySlice).Error; err != nil {
			return err
		}

		return nil
	})
}

func (_self sysPost) UserBindPostIds(db *gorm.DB, userId uint64, postIds []uint64) (err error) {

	if userId <= 0 {
		return errors.New("无用户信息")
	}
	if err = db.Delete(&entity.SysPostRelation{},
		"relation_id = ? and relation_type = ?", userId, consts.User).Error; err != nil {
		return
	}
	if len(postIds) > 0 {
		var slice []entity.SysPostRelation
		for _, roleId := range postIds {
			slice = append(slice, entity.SysPostRelation{
				PostId:       roleId,
				RelationId:   userId,
				RelationType: consts.User,
			})
		}
		err = db.Create(&slice).Error
	}

	return
}

// JoinUserPostDB
// select pr.relation_id as user_id,
// GROUP_CONCAT(p.name SEPARATOR ',') AS post_names,
// GROUP_CONCAT(p.id SEPARATOR ',') AS post_ids
// from sys_post_relation pr  left join sys_post p on
// and pr.relation_type = 'user' group by pr.relation_id
func (_self sysPost) JoinUserPostDB(db *gorm.DB) *gorm.DB {
	return db.Table(table_name.SysPostRelation+" as pr").
		Joins("left join "+table_name.SysPost+
			" as p on pr.post_id = p.id and pr.relation_type = ?", consts.User).
		Select("pr.relation_id as user_id",
			"GROUP_CONCAT(p.name SEPARATOR ',')  AS post_names",
			"GROUP_CONCAT(p.id SEPARATOR ',')  AS post_ids",
		).Group("pr.relation_id")
}
