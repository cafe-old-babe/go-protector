package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_result"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"go-protector/server/service"
)

var SysDictApi sysDict

type sysDict struct {
	base.Api
}

// region DictType

// DictTypePage 查询字段类型 分页查询
func (_self sysDict) DictTypePage(c *gin.Context) {
	var typeService service.SysDictType
	_self.MakeService(c, &typeService)
	var pageReq dto.DictTypePageReq
	if err := c.BindJSON(&pageReq); err != nil {
		typeService.Logger.Error("dictType page Error: %v", err)
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
		typeService.Logger.Error("dictType insert Error: %v", err)
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
		typeService.Logger.Error("dictType DictTypeDelete Error: %v", err)
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
		dataService.Logger.Error("dictData page  Error: %v", err)
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
		dataService.Logger.Error("dictData insert Error: %v", err)
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
		dataService.Logger.Error("dictData dictDataUpdateStatus Error: %v", err)
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
		dataService.Logger.Error("dictData DictDataDelete Error: %v", err)
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
