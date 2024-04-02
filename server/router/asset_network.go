package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initAssetNetwork(group *gin.RouterGroup) {
	routerGroup := group.Group("asset-network")
	{
		assetNetworkApi := api.AssetNetwork
		routerGroup.POST("/page", assetNetworkApi.Page)
		routerGroup.POST("/save", assetNetworkApi.Save)
		routerGroup.POST("/delete", assetNetworkApi.Delete)
	}
}
