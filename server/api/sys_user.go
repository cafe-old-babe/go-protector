package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/custom/c_logger"
	"go-protector/server/commons/custom/c_result"
	"go-protector/server/models/dto"
	"go-protector/server/service"
)

var UserApi user

type user struct{}

func (_self user) Login(c *gin.Context) {
	var loginDTO dto.Login
	if err := c.BindJSON(&loginDTO); err != nil {
		c_logger.GetLogger(c).Error("login Error: %v", err)
		c_result.Failure(c, nil, err.Error())
		return
	}
	var userService service.SysUser
	userService.MakeService(c)
	res := userService.DoLogin(loginDTO)
	c_result.Result(c, res)
}

func (_self user) Logout(c *gin.Context) {
	c_result.Success(c, nil)
}

func (_self user) SetStatus(c *gin.Context) {
	var updateDTO dto.SetStatus
	if err := c.BindJSON(&updateDTO); err != nil {
		c_logger.GetLogger(c).Error("SetStatus Error: %v", err)
		c_result.FailureErr(c, err)
		return
	}
	sysUserService := service.MakeSysUser(c)
	if err := sysUserService.SetStatus(&updateDTO); err != nil {
		c_result.FailureErr(c, err)
		return
	}
	c_result.Success(c, nil)

}
