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
		routerGroup := group.Group("asset-auth")
		{
			routerGroup.POST("/page", _assetAuth.Page)
			routerGroup.POST("/save", _assetAuth.Save)
			routerGroup.POST("/delete", _assetAuth.Delete)
			routerGroup.POST("/excel/import", _assetAuth.Import)
			routerGroup.POST("/excel/template", _assetAuth.ExportTemplate)
			routerGroup.POST("/excel/data", _assetAuth.ExportData)
		}
	})

}

var _assetAuth assetAuth

type assetAuth struct {
	base.Router
}

func (_self assetAuth) Page(c *gin.Context) {
	var pageReq dto.AssetAuthPageReq
	if err := c.ShouldBindJSON(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetAuthService service.AssetAuth
	_self.MakeService(c, &assetAuthService)
	c_result.Result(c, assetAuthService.Page(&pageReq))

}

func (_self assetAuth) Save(c *gin.Context) {
	var data entity.AssetAuth
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c_result.Err(c, err)
		return
	}
	var assetAuthService service.AssetAuth
	_self.MakeService(c, &assetAuthService)
	result := assetAuthService.SimpleSave(&data, func() error {
		return assetAuthService.SaveCheck(&data)
	})
	c_result.Result(c, result)
}

func (_self assetAuth) Delete(c *gin.Context) {
	var idsReq base.IdsReq
	err := c.ShouldBindJSON(&idsReq)
	if err != nil {
		c_result.Err(c, err)
		return
	}
	var assetAuthService service.AssetAuth
	_self.MakeService(c, &assetAuthService)
	idsReq.Value = &entity.AssetAuth{}
	result := assetAuthService.SimpleDelByIds(&idsReq)
	c_result.Result(c, result)
}

func (_self assetAuth) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c_result.Err(c, err)
		return
	}

	var assetAuthService service.AssetAuth
	_self.MakeService(c, &assetAuthService)

	if err = assetAuthService.Import(file); err != nil {
		c_result.Err(c, err)
		return
	}

	return
}

func (_self assetAuth) ExportTemplate(c *gin.Context) {

	var assetAuthService service.AssetAuth
	_self.MakeService(c, &assetAuthService)
	if err := assetAuthService.ExportTemplate(); err != nil {
		c_result.Err(c, err)
	}
	return
}

func (_self assetAuth) ExportData(c *gin.Context) {
	var pageReq dto.AssetAuthPageReq
	if err := c.ShouldBindJSON(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var assetAuthService service.AssetAuth
	_self.MakeService(c, &assetAuthService)
	if err := assetAuthService.ExportData(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}

	return
}
