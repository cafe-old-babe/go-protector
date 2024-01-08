package middleware

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/database"
	"go-protector/server/core/local"
)

func SetDB(ctx *gin.Context) {
	db := database.GetDB(ctx)
	ctx.Set(local.CtxKeyDB, db)
	ctx.Next()
}
