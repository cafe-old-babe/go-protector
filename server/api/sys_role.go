package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_result"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"go-protector/server/service"
	"strconv"
)

var SysRoleApi sysRole

type sysRole struct {
	base.Api
}

// Page 分页查询
func (_self sysRole) Page(c *gin.Context) {
	var req dto.SysRolePageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}

	var roleService service.SysRole
	_self.MakeService(c, &roleService)
	c_result.Result(c, roleService.Page(&req))
}

// Save 保存 insert or update
func (_self sysRole) Save(c *gin.Context) {

	var model entity.SysRole
	if err := c.ShouldBindJSON(&model); err != nil {
		c_result.Err(c, err)
		return
	}

	var roleService service.SysRole
	_self.MakeService(c, &roleService)

	result := roleService.SimpleSave(&model, func() error {
		return roleService.SaveCheck(&model)
	})
	c_result.Result(c, result)
}

// Delete 删除
func (_self sysRole) Delete(c *gin.Context) {
	var req dto.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var roleService service.SysRole
	_self.MakeService(c, &roleService)
	c_result.Result(c, roleService.Delete(&req))

}

// GetPermission 权限查询
func (_self sysRole) GetPermission(c *gin.Context) {

	roleIdStr := c.Param("roleId")
	var roleId uint64
	var err error
	if roleId, err = strconv.ParseUint(roleIdStr, 10, 64); err != nil {
		c_result.Err(c, c_error.ErrParamInvalid)
		return
	}
	var roleService service.SysRole
	_self.MakeService(c, &roleService)
	c_result.Result(c, roleService.GetPermission(roleId))
}

func (_self sysRole) SavePermission(c *gin.Context) {
	roleIdStr := c.Param("roleId")
	var roleId uint64
	var err error
	if roleId, err = strconv.ParseUint(roleIdStr, 10, 64); err != nil {
		c_result.Err(c, c_error.ErrParamInvalid)
		return
	}
	var req dto.IdsReq
	if err = c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var roleService service.SysRole
	_self.MakeService(c, &roleService)
	c_result.Result(c, roleService.SavePermission(roleId, req.GetIds()))

}

func (_self sysRole) SetStatus(c *gin.Context) {

	var roleId, status uint64
	var err error
	if roleId, err = strconv.ParseUint(c.Param("roleId"), 10, 64); err != nil {
		c_result.Err(c, c_error.ErrParamInvalid)
		return
	}

	if status, err = strconv.ParseUint(c.Param("status"), 10, 8); err != nil {
		c_result.Err(c, c_error.ErrParamInvalid)
		return
	}
	var roleService service.SysRole
	_self.MakeService(c, &roleService)
	c_result.Result(c, roleService.SetStatus(roleId, int8(status)))
}
