package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("menu")
		{
			routerGroup.POST("/list", _sysMenu.ListTree)
			routerGroup.POST("/save", _sysMenu.Save)
			routerGroup.POST("/delete", _sysMenu.Delete)
			//routerGroup.POST("/data/:id/:status", sysDictApi.DictDataUpdateStatus)

		}
	})
}

var _sysMenu sysMenu

type sysMenu struct {
}

// ListTree 树形列表
func (_self sysMenu) ListTree(c *gin.Context) {

	sysMenuService := service.MakeSysMenuService(c)
	res := sysMenuService.ListTree()
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
	var req base.IdsReq
	sysMenuService := service.MakeSysMenuService(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}

	res := sysMenuService.Delete(&req)
	c_result.Result(c, res)
}
