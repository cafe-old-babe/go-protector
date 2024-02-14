package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSysUser(group *gin.RouterGroup) {
	routerGroup := group.Group("user")
	{
		userApi := api.SysUserApi
		routerGroup.POST("login", userApi.Login)
		routerGroup.GET("info", userApi.UserInfo)
		routerGroup.GET("nav", userApi.Nav)
		routerGroup.POST("logout", userApi.Logout)
		routerGroup.POST("setStatus", userApi.SetStatus)
		routerGroup.POST("page", userApi.Page)
		deptGroup := routerGroup.Group("dept")
		{
			deptGroup.POST("tree", userApi.DeptTree)
			deptGroup.POST("save", userApi.DeptSave)
			deptGroup.POST("delete", userApi.DeptDelete)
		}
	}
}
