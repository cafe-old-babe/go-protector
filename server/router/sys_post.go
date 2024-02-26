package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSysPost(group *gin.RouterGroup) {
	routerGroup := group.Group("post")
	{
		sysPostApi := api.SysPostApi
		routerGroup.POST("/page", sysPostApi.Page)
		routerGroup.POST("/save", sysPostApi.Save)
		routerGroup.POST("/delete", sysPostApi.Delete)
	}
}
