package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("asset-group")
		{
			routerGroup.POST("/tree", _assetGroup.Tree)
			routerGroup.POST("/save", _assetGroup.Save)
			routerGroup.POST("/delete", _assetGroup.Delete)
		}
	})
}

var _assetGroup assetGroup

type assetGroup struct {
	base.Api
}

func (_self assetGroup) Tree(c *gin.Context) {
	var groupService service.AssetGroup
	_self.MakeService(c, &groupService)
	c_result.Result(c, groupService.Tree())
}

func (_self assetGroup) Save(c *gin.Context) {
	var model entity.AssetGroup
	// ShouldBindJSON http: 200  BindJSON: 400
	if err := c.ShouldBindJSON(&model); err != nil {
		c_result.Err(c, err)
		return
	}
	var groupService service.AssetGroup
	_self.MakeService(c, &groupService)
	res := groupService.SimpleSave(&model, func() error {
		return groupService.SaveCheck(&model)
	})
	c_result.Result(c, res)
}

func (_self assetGroup) Delete(c *gin.Context) {
	var req base.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var groupService service.AssetGroup
	_self.MakeService(c, &groupService)
	c_result.Result(c, groupService.DeleteByIds(&req))
}
