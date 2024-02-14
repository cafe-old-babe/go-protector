package c_result

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/custom/c_logger"
	"go-protector/server/core/custom/c_translator"
	"go-protector/server/models/dto"
	"net/http"
)

// Success 返回200
func Success(c *gin.Context, data any, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, dto.ResultSuccess(data, msg...))
}

// Failure 返回200
func Failure(c *gin.Context, data any, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, dto.ResultFailure(data, msg...))
}

func Err(c *gin.Context, err error) {
	err = c_translator.ConvertValidateErr(err)
	c_logger.GetLogger(c).Error("path: %s, err: %v", c.FullPath(), err)
	c.AbortWithStatusJSON(http.StatusOK, dto.ResultFailureErr(err))
}

// Custom 自定义返回
func Custom(c *gin.Context, code int, data any, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, dto.ResultCustom(code, data, msg...))
}

// Result 返回200
func Result(c *gin.Context, result *dto.Result) {
	c.AbortWithStatusJSON(http.StatusOK, result)
}
