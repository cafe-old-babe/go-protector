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
		routerGroup := group.Group("asset-account")
		{
			routerGroup.POST("/page", _assetAccount.Page)
			routerGroup.POST("/save", _assetAccount.Save)
			routerGroup.POST("/delete", _assetAccount.Delete)
		}
	})
}

var _assetAccount assetAccount

type assetAccount struct {
	base.Router
}

// Page 分页查询
func (_self assetAccount) Page(c *gin.Context) {
	var pageReq dto.AssetAccountPageReq
	if err := c.ShouldBindJSON(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var accountService service.AssetAccount
	_self.MakeService(c, &accountService)
	result := accountService.Page(&pageReq)
	c_result.Result(c, result)
}

// Save 接入从帐号
func (_self assetAccount) Save(c *gin.Context) {
	var model entity.AssetAccount
	if err := c.ShouldBindJSON(&model); err != nil {
		c_result.Err(c, err)
		return
	}
	var accountService service.AssetAccount
	_self.MakeService(c, &accountService)
	result := accountService.Save(&model)
	c_result.Result(c, result)
}

// Delete 删除
func (_self assetAccount) Delete(c *gin.Context) {
	var idsReq base.IdsReq
	if err := c.ShouldBindJSON(&idsReq); err != nil {
		c_result.Err(c, err)
		return
	}
	var accountService service.AssetAccount
	_self.MakeService(c, &accountService)
	//idsReq.Value = &entity.AssetAccount{}
	//result := accountService.SimpleDelByIds(&idsReq, func() (err error) {
	//	return accountService.CheckBatchDeleteByIds(idsReq.GetIds())
	//})
	result := accountService.Delete(&idsReq)
	c_result.Result(c, result)

}
