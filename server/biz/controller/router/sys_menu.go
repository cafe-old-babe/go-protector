package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/controller/api"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("menu")
		{
			sysDictApi := api.SysMenuApi
			routerGroup.POST("/list", sysDictApi.ListTree)
			routerGroup.POST("/save", sysDictApi.Save)
			routerGroup.POST("/delete", sysDictApi.Delete)
			//routerGroup.POST("/data/:id/:status", sysDictApi.DictDataUpdateStatus)

		}
	})
}
