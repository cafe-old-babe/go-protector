package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"go-protector/server/models/vo"
	"gorm.io/gorm"
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
			PID:        menu.PID,
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

func (_self *SysMenuService) Save(req *dto.SysMenuSaveReq) (result *dto.Result) {
	if req == nil {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
	}
	if err := binding.Validator.ValidateStruct(req); err != nil {
		return dto.ResultFailureErr(err)
	}
	if req.MenuType != 2 {
		// 如果不是按钮 组件名称必填
		if len(req.Component) <= 0 {
			return dto.ResultFailureErr(errors.New("组件名称必填"))
		}
	}
	model := &entity.SysMenu{
		ModelId:    entity.ModelId{ID: req.ID},
		Name:       req.Name,
		MenuType:   req.MenuType,
		PID:        req.PID,
		Permission: req.Permission,
		Hidden:     req.Hidden,
		Path:       req.Path,
		Component:  req.Component,
	}
	var res *gorm.DB
	if req.ID <= 0 {
		// 新增
		res = _self.DB.Create(model)
	} else {
		//更新
		res = _self.DB.Updates(model)
	}
	if res.Error != nil {
		return dto.ResultFailureErr(res.Error)
	}
	if res.RowsAffected <= 0 {
		return dto.ResultFailureMsg("保存失败")
	}
	return dto.ResultSuccess(model, "保存成功")
}

func (_self *SysMenuService) Delete(req *dto.IdsReq) *dto.Result {
	if err := _self.DB.Unscoped().Delete(&entity.SysMenu{}, req.GetIds()).Error; err != nil {
		return dto.ResultFailureErr(err)
	}
	// todo 删除关联关系
	return dto.ResultSuccessMsg("删除成功")
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
