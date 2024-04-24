package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_captcha"
	"go-protector/server/internal/custom/c_result"
)

var SystemApi system

type system struct {
	base.Api
}

// GenerateCaptcha 获取验证码
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
