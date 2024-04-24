package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
)

var AssetBasic assetBasic

type assetBasic struct {
	base.Api
}

func (_self assetBasic) Page(c *gin.Context) {
	var pageReq dto.AssetInfoPageReq
	if err := c.ShouldBindJSON(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetService service.AssetInfo
	_self.MakeService(c, &assetService)
	result := assetService.Page(&pageReq)
	c_result.Result(c, result)
}

func (_self assetBasic) Save(c *gin.Context) {
	var saveReq dto.AssetInfoSaveReq
	if err := c.ShouldBindJSON(&saveReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetService service.AssetInfo
	_self.MakeService(c, &assetService)
	result := assetService.Save(&saveReq)
	c_result.Result(c, result)

}

func (_self assetBasic) Delete(c *gin.Context) {
	var idsReq base.IdsReq
	if err := c.ShouldBindJSON(&idsReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetService service.AssetInfo
	_self.MakeService(c, &assetService)
	result := assetService.Delete(&idsReq)
	c_result.Result(c, result)
}

// Collectors 采集资产信息
func (_self assetBasic) Collectors(c *gin.Context) {
	collectorsType := c.Param("collectorsType")
	var idsReq base.IdsReq
	if err := c.ShouldBindJSON(&idsReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetService service.AssetInfo
	_self.MakeService(c, &assetService)
	result := assetService.Collectors(&idsReq, collectorsType)
	c_result.Result(c, result)
	return
}
