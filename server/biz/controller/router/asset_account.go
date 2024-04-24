package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/controller/api"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("asset-account")
		{
			assetAccountApi := api.AssetAccount
			routerGroup.POST("/page", assetAccountApi.Page)
			routerGroup.POST("/save", assetAccountApi.Save)
			routerGroup.POST("/delete", assetAccountApi.Delete)
		}
	})
}
