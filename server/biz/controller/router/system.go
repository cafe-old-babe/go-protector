package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/controller/api"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("system")
		{
			systemApi := api.SystemApi
			routerGroup.GET("/captcha", systemApi.GenerateCaptcha)
			routerGroup.GET("/routes", systemApi.Routes)
		}
	})
}
