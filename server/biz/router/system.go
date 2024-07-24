package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_captcha"
	"go-protector/server/internal/custom/c_result"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {

		group.GET("/ws/bus", _system.wsBus)
		routerGroup := group.Group("system")
		{
			routerGroup.GET("/captcha", _system.GenerateCaptcha)
			routerGroup.GET("/routes", _system.Routes)
		}
	})
}

var _system system

type system struct {
	base.Router
}

// GenerateCaptcha 获取验证码
// 3-4	后端-动态图片验证码
func (_self system) GenerateCaptcha(c *gin.Context) {

	id, b64s, err := c_captcha.Generate()
	if err != nil {
		c_result.Err(c, err)
		return
	}
	c_result.Result(c, base.ResultSuccess(map[string]string{
		"cid":  id,
		"b64s": b64s,
	}))

}

func (_self system) Routes(c *gin.Context) {
	c_result.Result(c, base.ResultSuccess(nil))
}

func (_self system) wsBus(c *gin.Context) {
	var sysStemService service.System
	_self.MakeService(c, &sysStemService)

	if err := sysStemService.Bus(); err != nil {
		c_result.Err(c, err)
		return
	}
}
