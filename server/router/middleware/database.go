package middleware

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/database"
	"go-protector/server/commons/local"
)

func SetDB(ctx *gin.Context) {
	db := database.GetDB(ctx)
	ctx.Set(local.CtxKeyDB, db)
	ctx.Next()
}
