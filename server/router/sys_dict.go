package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSysDict(group *gin.RouterGroup) {
	routerGroup := group.Group("dict")
	{
		sysDictApi := api.SysDictApi
		routerGroup.POST("/type", sysDictApi.DictTypePage)
		routerGroup.POST("/type/save", sysDictApi.DictTypeSave)
		routerGroup.POST("/type/delete", sysDictApi.DictTypeDelete)

		routerGroup.POST("/data", sysDictApi.DictDataPage)
		routerGroup.POST("/dataList/:dictType", sysDictApi.DictDataList)
		routerGroup.POST("/data/save", sysDictApi.DictDataSave)
		routerGroup.POST("/data/delete", sysDictApi.DictDataDelete)
		routerGroup.POST("/data/:id/:status", sysDictApi.DictDataUpdateStatus)

	}
}
