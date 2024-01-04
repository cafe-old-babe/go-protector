package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/custom/c_logger"
	"go-protector/server/commons/custom/result"
	"go-protector/server/models/dto"
	"go-protector/server/service"
)

var UserApi user

type user struct{}

func (_self user) Login(c *gin.Context) {
	var loginDTO dto.Login
	if err := c.BindJSON(&loginDTO); err != nil {
		c_logger.GetLogger(c).Error("login Error: %v", err)
		result.Failure(c, nil, err.Error())
		return
	}
	var userService service.SysUser
	userService.MakeService(c)
	res := userService.DoLogin(loginDTO)
	result.Result(c, res)
}

func (_self user) Logout(c *gin.Context) {
	result.Success(c, nil)
}
