package middleware

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/consts"
	"go-protector/server/core/database"
)

func SetDB(ctx *gin.Context) {
	db := database.GetDB(ctx)
	ctx.Set(consts.CtxKeyDB, db)
	ctx.Next()
}
