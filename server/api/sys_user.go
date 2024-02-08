package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_logger"
	"go-protector/server/core/custom/c_result"
	"go-protector/server/models/dto"
	"go-protector/server/service"
)

var SysUserApi sysUser

type sysUser struct {
	base.Api
}

func (_self sysUser) Login(c *gin.Context) {
	var loginDTO dto.Login
	if err := c.BindJSON(&loginDTO); err != nil {
		c_logger.GetLogger(c).Error("login Error: %v", err)
		c_result.Failure(c, nil, err.Error())
		return
	}
	var userService service.SysUser
	_self.MakeService(c, &userService)
	res := userService.DoLogin(loginDTO)
	c_result.Result(c, res)
}

func (_self sysUser) Logout(c *gin.Context) {
	c_result.Success(c, nil)
}

func (_self sysUser) SetStatus(c *gin.Context) {
	var updateDTO dto.SetStatus
	if err := c.BindJSON(&updateDTO); err != nil {
		c_logger.GetLogger(c).Error("SetStatus Error: %v", err)
		c_result.Err(c, err)
		return
	}
	var sysUserService service.SysUser
	_self.MakeService(c, &sysUserService)
	if err := sysUserService.SetStatus(&updateDTO); err != nil {
		c_result.Err(c, err)
		return
	}
	c_result.Success(c, nil)
}

func (_self sysUser) UserInfo(c *gin.Context) {
	var sysUserService service.SysUser
	_self.MakeService(c, &sysUserService)
	res := sysUserService.UserInfo()
	c_result.Result(c, res)
}

func (_self sysUser) Nav(c *gin.Context) {
	var sysUserService service.SysUser
	_self.MakeService(c, &sysUserService)
	res := sysUserService.Nav()
	c_result.Result(c, res)
}
