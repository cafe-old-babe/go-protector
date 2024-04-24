package c_result

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_logger"
	"go-protector/server/internal/custom/c_translator"
	"net/http"
)

// Success 返回200
func Success(c *gin.Context, data any, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, base.ResultSuccess(data, msg...))
}

// Failure 返回200
func Failure(c *gin.Context, data any, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, base.ResultFailure(data, msg...))
}

func Err(c *gin.Context, err error) {
	err = c_translator.ConvertValidateErr(err)
	c_logger.GetLogger(c).Error("path: %s, err: %v", c.FullPath(), err)
	c.AbortWithStatusJSON(http.StatusOK, base.ResultFailureErr(err))
}

// Custom 自定义返回
func Custom(c *gin.Context, code int, data any, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, base.ResultCustom(code, data, msg...))
}

// Result 返回200
func Result(c *gin.Context, result interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, result)
}
