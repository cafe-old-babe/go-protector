package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initLogin(group *gin.RouterGroup) {
	routerGroup := group.Group("user")
	{
		userApi := api.SysUserApi
		routerGroup.POST("login", userApi.Login)
		routerGroup.POST("logout", userApi.Logout)
		routerGroup.POST("setStatus", userApi.SetStatus)

	}
}
