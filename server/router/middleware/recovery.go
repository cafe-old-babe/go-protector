package middleware

import "github.com/gin-gonic/gin"

func Recovery(ctx *gin.Context) {
	if err := recover(); err != nil {
		if ctx.IsAborted() {

			ctx.Status(200)

		}
	}
}
