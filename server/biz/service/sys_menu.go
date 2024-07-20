package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/model/vo"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_error"
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

func (_self *SysMenuService) ListTree() (result *base.Result) {

	var menuSlice []entity.SysMenu

	if err := _self.GetDB().Find(&menuSlice).Error; err != nil {
		return base.ResultFailureErr(err)
	}
	menuMap := map[uint64][]vo.SysMenuVO{}

	for _, menu := range menuSlice {
		menuMap[menu.PID] = append(menuMap[menu.PID], vo.SysMenuVO{
			ID:         menu.ID,
			PID:        menu.PID,
			Name:       menu.Name,
			MenuType:   menu.MenuType,
			Permission: menu.Permission,
			Hidden:     menu.Hidden.Int16,
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
			Hidden:     menu.Hidden.Int16,
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

	return base.ResultSuccess(root)
}

func (_self *SysMenuService) Save(req *dto.SysMenuSaveReq) (result *base.Result) {
	if req == nil {
		return base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	if err := binding.Validator.ValidateStruct(req); err != nil {
		return base.ResultFailureErr(err)
	}
	if req.MenuType != 2 {
		// 如果不是按钮 组件名称必填
		if len(req.Component) <= 0 {
			return base.ResultFailureErr(errors.New("组件名称必填"))
		}
	}
	model := &entity.SysMenu{
		ModelId:    entity.ModelId{ID: req.ID},
		Name:       req.Name,
		MenuType:   req.MenuType,
		PID:        req.PID,
		Permission: req.Permission,
		Hidden: sql.NullInt16{
			Int16: req.Hidden,
			Valid: true,
		},
		Path:      req.Path,
		Component: req.Component,
	}
	var res *gorm.DB
	if req.ID <= 0 {
		// 新增
		res = _self.GetDB().Create(model)
	} else {
		//更新
		res = _self.GetDB().Updates(model)
	}
	if res.Error != nil {
		return base.ResultFailureErr(res.Error)
	}
	if res.RowsAffected <= 0 {
		return base.ResultFailureMsg("保存失败")
	}
	return base.ResultSuccess(model, "保存成功")
}

func (_self *SysMenuService) Delete(req *base.IdsReq) *base.Result {
	if err := _self.GetDB().Unscoped().Delete(&entity.SysMenu{}, req.GetIds()).Error; err != nil {
		return base.ResultFailureErr(err)
	}
	// todo 删除关联关系
	return base.ResultSuccessMsg("删除成功")
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
