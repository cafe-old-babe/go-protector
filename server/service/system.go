package service

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
)

type System struct {
	base.Service
}

func MakeSystem(c *gin.Context) *System {
	var self System
	self.Make(c)
	return &self
}
