package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/controller/api"
)

func init() {

	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("asset-gateway")
		{
			assetGatewayApi := api.AssetGateway
			routerGroup.POST("/page", assetGatewayApi.Page)
			routerGroup.POST("/save", assetGatewayApi.Save)
			routerGroup.POST("/delete", assetGatewayApi.Delete)
		}
	})

}
