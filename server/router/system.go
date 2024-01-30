package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSystem(group *gin.RouterGroup) {
	routerGroup := group.Group("system")
	{
		systemApi := api.SystemApi
		routerGroup.GET("/captcha", systemApi.GenerateCaptcha)
		routerGroup.GET("/routes", systemApi.Routes)

	}
}
