package base

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/custom/c_logger"
	"go-protector/server/core/database"
	"gorm.io/gorm"
)

type Service struct {
	Logger *c_logger.SelfLogger
	DB     *gorm.DB
	Ctx    *gin.Context
}

func (_self *Service) MakeService(ctx *gin.Context) {
	_self.DB = database.GetDB(ctx)
	_self.Logger = c_logger.GetLogger(ctx)
	_self.Ctx = ctx
}
