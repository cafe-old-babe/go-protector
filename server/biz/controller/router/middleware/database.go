package middleware

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/database"
)

func SetDB(ctx *gin.Context) {
	db := database.GetDB(ctx)
	ctx.Set(consts.CtxKeyDB, db)
	ctx.Next()
}
