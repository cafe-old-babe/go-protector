package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/logger"
	"go-protector/server/commons/result"
	"go-protector/server/models/dto"
	"go-protector/server/service"
)

var UserApi = new(User)

type User struct{}

func (_self User) Login(c *gin.Context) {
	var loginDTO dto.Login
	if err := c.BindJSON(&loginDTO); err != nil {
		logger.NewLogger(c).Error("login Error: %v", err)
		result.Failure(c, nil, err.Error())
		return
	}
	var userService service.SysUser
	userService.MakeService(c)
	res := userService.DoLogin(loginDTO)
	result.Result(c, res)
}

func (_self User) Logout(c *gin.Context) {
	result.Success(c, nil)
}
