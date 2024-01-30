package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSysDept(group *gin.RouterGroup) {
	routerGroup := group.Group("dept")
	{
		sysDeptApi := api.SysDeptApi
		routerGroup.POST("/page", sysDeptApi.Page)
		routerGroup.POST("/save", sysDeptApi.Save)
		routerGroup.POST("/delete", sysDeptApi.Delete)

	}
}
