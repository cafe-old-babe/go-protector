package service

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"go-protector/server/models/vo"
)

type SysMenuService struct {
	base.Service
}

func MakeSysMenuService(c *gin.Context) *SysMenuService {
	var self SysMenuService
	self.Make(c)
	return &self
}

func (_self *SysMenuService) List() (result *dto.Result) {

	var menuSlice []entity.SysMenu

	if err := _self.DB.Find(&menuSlice).Error; err != nil {
		return dto.ResultFailureErr(err)
	}
	menuMap := map[uint64][]vo.SysMenuVO{}

	for _, menu := range menuSlice {
		menuMap[menu.PID] = append(menuMap[menu.PID], vo.SysMenuVO{
			ID:         menu.ID,
			PID:        0,
			Name:       menu.Name,
			MenuType:   menu.MenuType,
			Permission: menu.Permission,
			Hidden:     menu.Hidden,
			Component:  menu.Component,
			Children:   []vo.SysMenuVO{},
		})
	}
	var menuVOSlice []vo.SysMenuVO
	for _, menu := range menuSlice {
		if menu.PID != 0 {
			continue
		}
		sysMenuVO := vo.SysMenuVO{
			ID:         menu.ID,
			PID:        0,
			Name:       menu.Name,
			MenuType:   menu.MenuType,
			Permission: menu.Permission,
			Hidden:     menu.Hidden,
			Component:  menu.Component,
			Children:   []vo.SysMenuVO{},
		}

		sysMenuVO.Children = generateChildren(&sysMenuVO, menuMap)
		menuVOSlice = append(menuVOSlice, sysMenuVO)
	}
	root := vo.SysMenuVO{
		ID:       0,
		Name:     "根节点",
		Children: menuVOSlice,
	}

	return dto.ResultSuccess(root)
}

func generateChildren(menuVO *vo.SysMenuVO, menuMap map[uint64][]vo.SysMenuVO) (children []vo.SysMenuVO) {

	children = menuMap[menuVO.ID]

	switch menuVO.MenuType {
	case 0:
		menuVO.MenuTypeName = "目录"
	case 1:
		menuVO.MenuTypeName = "菜单"
	case 2:
		menuVO.MenuTypeName = "按钮"
		return children
	}
	if len(children) <= 0 {
		return nil
	}
	for i := range children {
		children[i].Children = generateChildren(&children[i], menuMap)
	}

	return children

}
