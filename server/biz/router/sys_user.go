package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_logger"
	"go-protector/server/internal/custom/c_result"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("user")
		{
			routerGroup.POST("login", _sysUser.Login)
			routerGroup.GET("info", _sysUser.UserInfo)
			routerGroup.GET("nav", _sysUser.Nav)
			routerGroup.POST("logout", _sysUser.Logout)
			routerGroup.POST("setStatus", _sysUser.SetStatus)
			routerGroup.POST("page", _sysUser.Page)
			routerGroup.POST("save", _sysUser.Save)
			routerGroup.POST("delete", _sysUser.Delete)
			deptGroup := routerGroup.Group("dept")
			{
				deptGroup.POST("tree", _sysUser.DeptTree)
				deptGroup.POST("save", _sysUser.DeptSave)
				deptGroup.POST("delete", _sysUser.DeptDelete)
			}
		}
	})
}

var _sysUser sysUser

type sysUser struct {
	base.Api
}

func (_self sysUser) Login(c *gin.Context) {
	var loginDTO dto.LoginDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c_logger.GetLogger(c).Error("login Error: %v", err)
		c_result.Failure(c, nil, err.Error())
		return
	}
	var userService service.SysUser
	_self.MakeService(c, &userService)
	res := userService.DoLogin(loginDTO)
	c_result.Result(c, res)
}

func (_self sysUser) Logout(c *gin.Context) {
	c_result.Success(c, nil)
}

func (_self sysUser) SetStatus(c *gin.Context) {
	var updateDTO dto.SetStatus
	if err := c.BindJSON(&updateDTO); err != nil {
		c_logger.GetLogger(c).Error("SetStatus Error: %v", err)
		c_result.Err(c, err)
		return
	}
	var sysUserService service.SysUser
	_self.MakeService(c, &sysUserService)
	if err := sysUserService.SetStatus(&updateDTO); err != nil {
		c_result.Err(c, err)
		return
	}
	c_result.Success(c, nil)
}

func (_self sysUser) UserInfo(c *gin.Context) {
	var sysUserService service.SysUser
	_self.MakeService(c, &sysUserService)
	res := sysUserService.UserInfo()
	c_result.Result(c, res)
}

func (_self sysUser) Nav(c *gin.Context) {
	var sysUserService service.SysUser
	_self.MakeService(c, &sysUserService)
	res := sysUserService.Nav()
	c_result.Result(c, res)
}

// Page 用户分页查询
func (_self sysUser) Page(c *gin.Context) {
	var sysUserService service.SysUser
	_self.MakeService(c, &sysUserService)
	var userPageReq dto.UserPageReq
	if err := c.ShouldBindJSON(&userPageReq); err != nil {
		c_result.Err(c, err)
		return
	}
	res := sysUserService.Page(&userPageReq)
	c_result.Result(c, res)
}

// DeptTree 部门树
func (_self sysUser) DeptTree(c *gin.Context) {

	var sysDeptService service.SysDept
	_self.MakeService(c, &sysDeptService)
	res := sysDeptService.DeptTree()
	c_result.Result(c, res)

}

// DeptDelete 删除
func (_self sysUser) DeptDelete(c *gin.Context) {
	var sysDeptService service.SysDept
	var ids base.IdsReq
	_self.MakeService(c, &sysDeptService)
	if err := c.ShouldBindJSON(&ids); err != nil {
		c_result.Err(c, err)
		return
	}
	ids.Unscoped = true
	ids.Value = &entity.SysDept{}
	res := sysDeptService.SimpleDelByIds(&ids, func() error {
		if len(ids.GetIds()) <= 0 {
			return c_error.ErrDeleteFailure
		}
		return nil
	})
	c_result.Result(c, res)
}

// DeptSave 保存
func (_self sysUser) DeptSave(c *gin.Context) {
	var sysDeptService service.SysDept
	var model entity.SysDept
	// ShouldBindJSON http: 200  BindJSON: 400
	if err := c.ShouldBindJSON(&model); err != nil {
		c_result.Err(c, err)
		return
	}
	_self.MakeService(c, &sysDeptService)
	res := sysDeptService.SimpleSave(&model, func() error {
		return sysDeptService.SaveCheck(&model)
	})
	c_result.Result(c, res)
}

func (_self sysUser) Save(c *gin.Context) {

	var req dto.UserSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var sysUserService service.SysUser
	_self.MakeService(c, &sysUserService)
	c_result.Result(c, sysUserService.Save(&req))

}

func (_self sysUser) Delete(c *gin.Context) {
	var req base.IdsReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var userService service.SysUser
	_self.MakeService(c, &userService)
	c_result.Result(c, userService.DeleteByIds(&req))
}
