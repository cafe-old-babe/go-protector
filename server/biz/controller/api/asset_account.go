package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
)

var AssetAccount assetAccount

type assetAccount struct {
	base.Api
}

func (_self assetAccount) Page(c *gin.Context) {
	var pageReq dto.AssetAccountPageReq
	if err := c.ShouldBindJSON(pageReq); err != nil {
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

}
func (_self assetAccount) Delete(c *gin.Context) {

}
