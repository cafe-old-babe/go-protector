package base

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/custom/c_logger"
	"go-protector/server/core/database"
	"gorm.io/gorm"
)

type IService interface {
	Make(ctx *gin.Context)
	MakeService(service ...IService)
}

type Service struct {
	Logger  *c_logger.SelfLogger
	DB      *gorm.DB
	Context *gin.Context
}

func (_self *Service) Make(c *gin.Context) {
	_self.DB = database.GetDB(c)
	_self.Logger = c_logger.GetLogger(c)
	_self.Context = c
}

func (_self *Service) MakeService(service ...IService) {
	if len(service) <= 0 {
		return
	}
	for i := range service {
		service[i].Make(_self.Context)
	}
}
