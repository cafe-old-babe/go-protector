package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
	"go-protector/server/core/custom/c_result"
	"go-protector/server/core/custom/c_type"
	"go-protector/server/service"
)

var SysLoginPolicyApi sysLoginPolicy

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
