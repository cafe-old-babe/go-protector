package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSysRole(group *gin.RouterGroup) {
	routerGroup := group.Group("role")
	{
		sysRoleApi := api.SysRoleApi
		routerGroup.POST("/page", sysRoleApi.Page)
		routerGroup.POST("/save", sysRoleApi.Save)
		routerGroup.POST("/getPermission/:roleId", sysRoleApi.GetPermission)
		routerGroup.POST("/savePermission/:roleId", sysRoleApi.SavePermission)
		routerGroup.POST("/setStatus/:roleId/:status", sysRoleApi.SetStatus)
		routerGroup.POST("/delete", sysRoleApi.Delete)
		//routerGroup.POST("/data/:id/:status", sysDictApi.DictDataUpdateStatus)

	}
}
