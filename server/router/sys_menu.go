package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSysMenu(group *gin.RouterGroup) {
	routerGroup := group.Group("menu")
	{
		sysDictApi := api.SysMenuApi
		routerGroup.POST("/list", sysDictApi.List)
		routerGroup.POST("/save", sysDictApi.Save)
		routerGroup.POST("/delete", sysDictApi.Delete)
		//routerGroup.POST("/data/:id/:status", sysDictApi.DictDataUpdateStatus)

	}
}
