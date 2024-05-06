package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
	"go-protector/server/internal/custom/c_type"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("sys-login-policy")
		{
			routerGroup.POST("/info", _sysLoginPolicy.Info)
			routerGroup.POST("/save", _sysLoginPolicy.Save)

		}
	})
}

var _sysLoginPolicy sysLoginPolicy

type sysLoginPolicy struct {
	base.Api
}

func (_self sysLoginPolicy) Info(c *gin.Context) {
	var policyService service.SysLoginPolicy
	_self.MakeService(c, &policyService)
	c_result.Result(c, policyService.Info())
}

func (_self sysLoginPolicy) Save(c *gin.Context) {
	var policyService service.SysLoginPolicy
	_self.MakeService(c, &policyService)
	var param map[c_type.LoginPolicyCode]map[string]interface{}
	if err := c.ShouldBindJSON(&param); err != nil {
		c_result.Err(c, err)
		return
	}
	c_result.Result(c, policyService.Save(param))
}
