package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/controller/api"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("asset-group")
		{
			assetGroupApi := api.AssetGroup
			routerGroup.POST("/tree", assetGroupApi.Tree)
			routerGroup.POST("/save", assetGroupApi.Save)
			routerGroup.POST("/delete", assetGroupApi.Delete)
		}
	})
}
