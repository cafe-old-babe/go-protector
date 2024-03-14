package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSysLoginPolicy(group *gin.RouterGroup) {
	routerGroup := group.Group("sys-login-policy")
	{
		sysLoginPolicyApi := api.SysLoginPolicyApi
		routerGroup.POST("/info", sysLoginPolicyApi.Info)
		routerGroup.POST("/save", sysLoginPolicyApi.Save)

	}
}
