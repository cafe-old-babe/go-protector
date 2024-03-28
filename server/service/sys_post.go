package service

import (
	"errors"
	"fmt"
	"go-protector/server/core/base"
	"go-protector/server/core/consts"
	"go-protector/server/core/consts/table_name"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/scope"
	"go-protector/server/dao"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"go-protector/server/models/vo"
	"gorm.io/gorm"
	"strings"
)

type SysPost struct {
	base.Service
}

func (_self *SysPost) Page(req *dto.SysPostPageReq) *base.Result {
	var slice []vo.PostPage
	var count int64

	/*
		select p.id,p.name,p.code, GROUP_CONCAT(r.id) r_ids,
		       GROUP_CONCAT(d.dept_name) as dept_names,GROUP_CONCAT(d.id) as dept_ids
		from sys_post p
		left join  sys_post_relation r on p.id = r.post_id and r.relation_type='dept'
		left join (
		    WITH RECURSIVE sys_dept_path AS (
		        SELECT id, dept_name  FROM sys_dept WHERE p_id = 0 -- 根节点为空
		        UNION ALL
		        SELECT d.id, CONCAT(d.dept_name, '/',dp.dept_name ) as dept_path
		        FROM sys_dept d INNER JOIN sys_dept_path dp ON d.p_id = dp.id
		    ) SELECT * FROM sys_dept_path
		) d on r.relation_id = d.id group by p.id;
	*/
	tx := _self.DB.Table(table_name.SysPost + " as p")
	builder := strings.Builder{}
	builder.WriteString("WITH RECURSIVE sys_dept_path AS ( ")
	builder.WriteString("SELECT id, dept_name  FROM sys_dept WHERE p_id = 0 ")
	builder.WriteString("UNION ALL ")
	builder.WriteString("SELECT d.id, CONCAT(d.dept_name, '/',dp.dept_name ) as dept_path ")
	builder.WriteString("FROM sys_dept d INNER JOIN sys_dept_path dp ON d.p_id = dp.id ")
	builder.WriteString(") SELECT * FROM sys_dept_path ")
	rawTx := _self.DB.Raw(builder.String())
	tx = tx.
		Select([]string{"p.id", "p.name", "p.code", "GROUP_CONCAT(r.id) r_ids", "GROUP_CONCAT(d.id) as dept_ids"}).
		Scopes(
			scope.Paginate(req.GetPagination()),
			scope.Like("name", req.Name),
			scope.Like("LOWER(code)", strings.ToLower(req.Code)),
		)
	joinRelationSqlFormat := " %s join " + table_name.SysPostRelation + " as r on p.id = r.post_id and r.relation_type= ? %s"
	if len(req.DeptIds) <= 0 {
		tx = tx.Joins(fmt.Sprintf(joinRelationSqlFormat, "left", ""), consts.Dept)
	} else {
		tx = tx.Joins(fmt.Sprintf(joinRelationSqlFormat, "inner", "and r.relation_id in (?)"), consts.Dept, req.DeptIds)
	}
	if err := tx.Joins("left join (?) as d on r.relation_id = d.id", rawTx).Group("p.id").
		Find(&slice).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return base.ResultFailureErr(err)
	}
	return base.ResultPage(slice, req.GetPagination(), count)
}

func (_self *SysPost) CheckSave(req *dto.SysPostSaveReq) (err error) {
	var count int64
	if err = _self.DB.Table(table_name.SysPost).
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
		err = errors.New("岗位名称和岗位代码不能重复")
	}

	return
}

func (_self *SysPost) DeleteByIds(req *base.IdsReq) *base.Result {

	if err := dao.SysPost.DeleteByPostId(_self.DB, req.GetIds()); err != nil {
		return base.ResultFailureErr(err)
	}
	return base.ResultSuccessMsg("删除成功")

}

func (_self *SysPost) Save(req *dto.SysPostSaveReq) *base.Result {

	if err := dao.SysPost.Save(_self.DB, req); err != nil {
		return base.ResultFailureErr(err)
	}
	return base.ResultSuccessMsg("保存成功")
}

func (_self *SysPost) ListByDeptId(deptId uint64) *base.Result {
	if deptId <= 0 {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	var postSlice []entity.SysPost
	err := _self.DB.Find(&postSlice, "id in (?)",
		_self.DB.Select("post_id").
			Model(&entity.SysPostRelation{}).
			Where("relation_type = ? and relation_id = ?", consts.Dept, deptId)).Error
	if err != nil {
		return base.ResultFailureErr(err)
	}
	return base.ResultSuccess(postSlice)

}
