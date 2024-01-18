package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/custom/c_captcha"
	"go-protector/server/core/custom/c_result"
	"go-protector/server/models/dto"
)

var SystemApi system

type system struct{}

// GenerateCaptcha 获取验证码
func (_self system) GenerateCaptcha(c *gin.Context) {
	id, b64s, err := c_captcha.Generate()
	if err != nil {
		c_result.FailureErr(c, err)
		return
	}
	c_result.Result(c, dto.ResultSuccess(map[string]string{
		"cid":  id,
		"b64s": b64s,
	}))

}

func (_self system) Routes(c *gin.Context) {
	c_result.Result(c, dto.ResultSuccess(nil))
}
