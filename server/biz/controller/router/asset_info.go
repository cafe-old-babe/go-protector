package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/controller/api"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("asset-info")
		{
			assetBasic := api.AssetBasic
			routerGroup.POST("/page", assetBasic.Page)
			routerGroup.POST("/save", assetBasic.Save)
			routerGroup.POST("/collectors/:collectorsType", assetBasic.Collectors)
			routerGroup.POST("/delete", assetBasic.Delete)
		}
	})
}
