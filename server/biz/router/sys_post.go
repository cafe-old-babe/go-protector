package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_result"
	"strconv"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("post")
		{
			routerGroup.POST("/page", _sysPost.Page)
			// 4-20	后端-部门与岗位管理-岗位接口增删改接口开发（GORM-DeletedAt注意事项）
			routerGroup.POST("/list/:deptId", _sysPost.List)
			routerGroup.POST("/save", _sysPost.Save)
			routerGroup.POST("/delete", _sysPost.Delete)
		}
	})
}

var _sysPost sysPost

type sysPost struct {
	base.Router
}

// Page 分页
func (_self sysPost) Page(c *gin.Context) {
	var req dto.SysPostPageReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var postService service.SysPost
	_self.MakeService(c, &postService)
	c_result.Result(c, postService.Page(&req))

}

// Save 保存
func (_self sysPost) Save(c *gin.Context) {
	var req dto.SysPostSaveReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var postService service.SysPost
	_self.MakeService(c, &postService)
	c_result.Result(c, postService.Save(&req))
}

// Delete 删除
func (_self sysPost) Delete(c *gin.Context) {
	var req base.IdsReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var postService service.SysPost
	req.Value = &entity.SysPost{}
	_self.MakeService(c, &postService)
	c_result.Result(c, postService.DeleteByIds(&req))
}

func (_self sysPost) List(c *gin.Context) {
	deptIdStr := c.Param("deptId")
	deptId, err := strconv.ParseUint(deptIdStr, 10, 64)
	if err != nil {
		c_result.Err(c, err)
		return
	}
	if deptId <= 0 {
		c_result.Err(c, c_error.ErrParamInvalid)
		return
	}
	var postService service.SysPost
	_self.MakeService(c, &postService)
	c_result.Result(c, postService.ListByDeptId(deptId))
}
