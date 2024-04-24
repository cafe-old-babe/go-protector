package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/controller/api"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("sys-login-policy")
		{
			sysLoginPolicyApi := api.SysLoginPolicyApi
			routerGroup.POST("/info", sysLoginPolicyApi.Info)
			routerGroup.POST("/save", sysLoginPolicyApi.Save)

		}
	})
}
