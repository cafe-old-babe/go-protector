package result

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/models/dto"
	"net/http"
)

// Success 返回200
func Success(c *gin.Context, data any, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, dto.ResultSuccess(data, msg...))
}

// Failure 返回400
func Failure(c *gin.Context, data any, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, dto.ResultFailure(data, msg...))
}

// Custom 返回400
func Custom(c *gin.Context, code int, data any, msg ...string) {
	c.AbortWithStatusJSON(http.StatusOK, dto.ResultCustom(code, data, msg...))
}

// Result 返回200
func Result(c *gin.Context, result *dto.Result) {
	c.AbortWithStatusJSON(http.StatusOK, result)
}
