package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
)

func init() {

	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("asset-gateway")
		{
			routerGroup.POST("/page", _assetGateway.Page)
			routerGroup.POST("/list", _assetGateway.List)
			routerGroup.POST("/save", _assetGateway.Save)
			routerGroup.POST("/delete", _assetGateway.Delete)
		}
	})

}

var _assetGateway assetGateway

type assetGateway struct {
	base.Router
}

func (_self assetGateway) Page(c *gin.Context) {
	var pageReq dto.AssetGatewayPageReq
	if err := c.ShouldBindJSON(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetGatewayService service.AssetGateway
	_self.MakeService(c, &assetGatewayService)

	c_result.Result(c, assetGatewayService.Page(&pageReq))

}

func (_self assetGateway) Save(c *gin.Context) {
	var model entity.AssetGateway
	if err := c.ShouldBindJSON(&model); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetGatewayService service.AssetGateway
	_self.MakeService(c, &assetGatewayService)
	result := assetGatewayService.SimpleSave(&model, func() error {
		return assetGatewayService.Check(&model)
	})
	c_result.Result(c, result)

}

func (_self assetGateway) Delete(c *gin.Context) {
	var idsReq base.IdsReq
	if err := c.ShouldBindJSON(&idsReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetGatewayService service.AssetGateway
	_self.MakeService(c, &assetGatewayService)
	idsReq.Value = &entity.AssetGateway{}
	c_result.Result(c, assetGatewayService.SimpleDelByIds(&idsReq, func() error {
		// todo 检查使用情况
		return nil
	}))
	return

}

func (_self assetGateway) List(c *gin.Context) {
	var assetGatewayService service.AssetGateway
	_self.MakeService(c, &assetGatewayService)
	c_result.Result(c, assetGatewayService.List())
	return

}
