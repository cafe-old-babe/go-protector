package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_result"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"go-protector/server/service"
)

var AssetNetwork assetNetwork

type assetNetwork struct {
	base.Api
}

func (_self assetNetwork) Page(c *gin.Context) {
	var pageReq dto.AssetNetworkPageReq
	if err := c.ShouldBindJSON(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetNetworkService service.AssetNetwork
	_self.MakeService(c, &assetNetworkService)

	c_result.Result(c, assetNetworkService.Page(&pageReq))

}

func (_self assetNetwork) Save(c *gin.Context) {
	var model entity.AssetNetwork
	if err := c.ShouldBindJSON(&model); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetNetworkService service.AssetNetwork
	_self.MakeService(c, &assetNetworkService)
	result := assetNetworkService.SimpleSave(&model, func() error {
		return assetNetworkService.Check(&model)
	})
	c_result.Result(c, result)

}

func (_self assetNetwork) Delete(c *gin.Context) {
	var idsReq base.IdsReq
	if err := c.ShouldBindJSON(&idsReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetNetworkService service.AssetNetwork
	_self.MakeService(c, &assetNetworkService)
	idsReq.Value = &entity.AssetNetwork{}
	c_result.Result(c, assetNetworkService.SimpleDelByIds(&idsReq, func() error {
		// todo 检查使用情况
		return nil
	}))
	return

}
