package middleware

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/database"
	"go-protector/server/commons/local"
)

func SetDB(ctx *gin.Context) {
	ctx.Set(local.CtxKeyDB, database.GetDB())
	ctx.Next()
}
