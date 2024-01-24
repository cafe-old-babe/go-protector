package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSysDict(group *gin.RouterGroup) {
	sysDict := group.Group("dict")
	{
		sysDictApi := api.SysDictApi
		sysDict.POST("/type", sysDictApi.DictTypePage)
		sysDict.POST("/type/save", sysDictApi.DictTypeSave)
		sysDict.POST("/type/delete", sysDictApi.DictTypeDelete)

		sysDict.POST("/data", sysDictApi.DictDataPage)
		sysDict.POST("/data/save", sysDictApi.DictDataSave)
		sysDict.POST("/data/delete", sysDictApi.DictDataDelete)
		sysDict.POST("/data/:id/:status", sysDictApi.DictDataUpdateStatus)

	}
}
