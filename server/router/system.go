package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func InitSystem(group *gin.RouterGroup) {
	system := group.Group("system")
	{
		systemApi := api.SystemApi
		system.GET("/captcha", systemApi.GenerateCaptcha)

	}
}
