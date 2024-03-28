package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_result"
	"go-protector/server/models/dto"
	"go-protector/server/service"
	"strconv"
)

var SysPostApi sysPost

type sysPost struct {
	base.Api
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
	req.Value = &service.SysPost{}
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
