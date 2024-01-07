package api

import (
	"github.com/gin-gonic/gin"
)

var SystemApi system

type system struct{}

func (s system) GenerateCaptcha(c *gin.Context) {

}
