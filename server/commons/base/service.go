package base

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/database"
	"go-protector/server/commons/logger"
	"gorm.io/gorm"
)

type Service struct {
	Logger *logger.CustomLogger
	DB     *gorm.DB
	Ctx    *gin.Context
}

func (_self *Service) MakeService(ctx *gin.Context) {
	_self.DB = database.GetDB()
	_self.Logger = logger.NewLogger(ctx)
	_self.Ctx = ctx
}
