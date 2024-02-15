package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/custom/c_result"
	"go-protector/server/models/dto"
	"go-protector/server/service"
)

var SysMenuApi sysMenu

type sysMenu struct {
}

// List 树形列表
func (_self sysMenu) List(c *gin.Context) {

	sysMenuService := service.MakeSysMenuService(c)
	res := sysMenuService.List()
	c_result.Result(c, res)

}

// Save 保存
func (_self sysMenu) Save(c *gin.Context) {
	sysMenuService := service.MakeSysMenuService(c)
	var req dto.SysMenuSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	res := sysMenuService.Save(&req)
	c_result.Result(c, res)
}

// Delete 删除
func (_self sysMenu) Delete(c *gin.Context) {
	var req dto.IdsReq
	sysMenuService := service.MakeSysMenuService(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}

	res := sysMenuService.Delete(&req)
	c_result.Result(c, res)
}
