package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func initSysDict(group *gin.RouterGroup) {
	sysDict := group.Group("dict")
	{
		sysDictApi := api.SysDictApi
		sysDict.GET("/type", sysDictApi.DictTypePage)
		sysDict.POST("/type/insert", sysDictApi.DictTypeInsert)
		sysDict.POST("/type/update", sysDictApi.DictTypeUpdate)
		sysDict.POST("/type/delete", sysDictApi.DictTypeDelete)

		sysDict.GET("/data", sysDictApi.DictDataPage)
		sysDict.POST("/data/insert", sysDictApi.DictDataInsert)
		sysDict.POST("/data/update", sysDictApi.DictDataUpdate)
		sysDict.POST("/data/delete", sysDictApi.DictDataDelete)
		sysDict.POST("/data/:id/:status", sysDictApi.DictDataUpdateStatus)

	}
}
