package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
	"strings"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("asset-info")
		{
			//6-8	【实战】资源、特权从帐号管理接口实战（GORM-Belongs To、Joins预加载，GO-使用反射+泛型+断言将切片中对象的某一列转为slice）
			routerGroup.POST("/page", _assetBasic.Page)
			routerGroup.POST("/save", _assetBasic.Save)
			routerGroup.POST("/delete", _assetBasic.Delete)

			// 6-15	【实战】采集资源从账号-掌握defer注意事项
			routerGroup.POST("/collectors/:collType", _assetBasic.Collectors)
			// 6-14	【实战】拨测资源从帐号-掌握策略模式；掌握GORM Has one、Preload预加载
			routerGroup.POST("/dial/:dialType", _assetBasic.Dial)
			routerGroup.POST("/auth/page", _assetBasic.Page)
		}
	})
}

var _assetBasic assetBasic

type assetBasic struct {
	base.Router
}

func (_self assetBasic) Page(c *gin.Context) {
	var pageReq dto.AssetInfoPageReq
	if err := c.ShouldBindJSON(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetService service.AssetInfo
	_self.MakeService(c, &assetService)
	// 8-2	【实战】我的资源-复用资源分页查询接口，查询已授权的资源
	pageReq.Auth = strings.HasSuffix(c.Request.RequestURI, "/auth/page")
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
	collType := c.Param("collType")
	var idsReq base.IdsReq
	if err := c.ShouldBindJSON(&idsReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetService service.AssetInfo
	_self.MakeService(c, &assetService)
	result := assetService.Collectors(&idsReq, collType)
	c_result.Result(c, result)
	return
}

// Dial 拨测账号
func (_self assetBasic) Dial(c *gin.Context) {
	dialType := c.Param("dialType")

	var idsReq base.IdsReq
	if err := c.ShouldBindJSON(&idsReq); err != nil {
		c_result.Err(c, err)
		return
	}

	var assetService service.AssetInfo
	_self.MakeService(c, &assetService)
	result := assetService.Dial(&idsReq, dialType)
	c_result.Result(c, result)
	return

}
