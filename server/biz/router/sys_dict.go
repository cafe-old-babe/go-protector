package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_result"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("dict")
		{
			routerGroup.POST("/type", _sysDict.DictTypePage)
			routerGroup.POST("/type/save", _sysDict.DictTypeSave)
			routerGroup.POST("/type/delete", _sysDict.DictTypeDelete)

			routerGroup.POST("/data", _sysDict.DictDataPage)
			routerGroup.POST("/dataList/:dictType", _sysDict.DictDataList)
			routerGroup.POST("/data/save", _sysDict.DictDataSave)
			routerGroup.POST("/data/delete", _sysDict.DictDataDelete)
			routerGroup.POST("/data/:id/:status", _sysDict.DictDataUpdateStatus)

		}
	})
}

var _sysDict sysDict

type sysDict struct {
	base.Router
}

// region DictType

// DictTypePage 查询字段类型 分页查询
func (_self sysDict) DictTypePage(c *gin.Context) {
	var typeService service.SysDictType
	_self.MakeService(c, &typeService)
	var pageReq dto.DictTypePageReq
	if err := c.BindJSON(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}

	res := typeService.Page(&pageReq)
	c_result.Result(c, res)
}

// DictTypeSave 类型保存
func (_self sysDict) DictTypeSave(c *gin.Context) {
	var typeService service.SysDictType
	_self.MakeService(c, &typeService)
	var model entity.SysDictType
	if err := c.BindJSON(&model); err != nil {
		c_result.Err(c, err)
		return
	}

	c_result.Result(c, typeService.Save(&model))
}

// DictTypeDelete 删除类型
func (_self sysDict) DictTypeDelete(c *gin.Context) {
	var typeService service.SysDictType
	_self.MakeService(c, &typeService)
	var req base.IdsReq
	if err := c.BindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	c_result.Result(c, typeService.Delete(&req))
}

// endregion DictType

//region DictData

// DictDataPage 查询字典数据
func (_self sysDict) DictDataPage(c *gin.Context) {
	var dataService service.SysDictData
	_self.MakeService(c, &dataService)
	var pageReq dto.DictDataPageReq
	if err := c.BindJSON(&pageReq); err != nil {
		c_result.Err(c, err)
		return
	}

	res := dataService.Page(&pageReq)
	c_result.Result(c, res)

}

// DictDataSave 数据保存
func (_self sysDict) DictDataSave(c *gin.Context) {
	var dataService service.SysDictData
	_self.MakeService(c, &dataService)
	var model entity.SysDictData
	if err := c.BindJSON(&model); err != nil {
		c_result.Err(c, err)
		return
	}
	c_result.Result(c, dataService.Save(&model))
}

// DictDataUpdateStatus 更新数据状态
func (_self sysDict) DictDataUpdateStatus(c *gin.Context) {
	var dataService service.SysDictData
	_self.MakeService(c, &dataService)
	var updateReq dto.DictDataUpdateStatusReq
	if err := c.BindUri(&updateReq); err != nil {
		c_result.Err(c, err)
		return
	}
	c_result.Result(c, dataService.UpdateStatus(&updateReq))
}

func (_self sysDict) DictDataDelete(c *gin.Context) {
	var dataService service.SysDictData
	_self.MakeService(c, &dataService)
	var req base.IdsReq
	if err := c.BindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	c_result.Result(c, dataService.Delete(&req))
}

func (_self sysDict) DictDataList(c *gin.Context) {
	var dataService service.SysDictData
	_self.MakeService(c, &dataService)
	dictType := c.Param("dictType")
	if len(dictType) <= 0 {
		c_result.Err(c, c_error.ErrParamInvalid)
		return
	}
	c_result.Result(c, dataService.DictDataList(&dictType))
}

// endregion region
