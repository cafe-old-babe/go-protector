package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/controller/api"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("post")
		{
			sysPostApi := api.SysPostApi
			routerGroup.POST("/page", sysPostApi.Page)
			routerGroup.POST("/list/:deptId", sysPostApi.List)
			routerGroup.POST("/save", sysPostApi.Save)
			routerGroup.POST("/delete", sysPostApi.Delete)
		}
	})
}
