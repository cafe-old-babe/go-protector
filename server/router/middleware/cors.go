package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 设置跨域
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		//设置允许的源。这里设置为 "*",表示允许所有源的请求。
		ctx.Header("Access-Control-Allow-Origin", "*")
		//设置允许的请求方法。这里设置了常见的 HTTP 方法,包括 POST、GET、OPTIONS、PUT、DELETE 和 UPDATE。
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		//设置允许的请求头。这里设置为 "*",表示允许所有请求头的请求。
		ctx.Header("Access-Control-Allow-Headers", "*")
		//设置响应头需要暴露给客户端的头部字段。这里设置了几个常用的响应头字段,包括 Content-Length、Access-Control-Allow-Origin、Access-Control-Allow-Headers、Cache-Control、Content-Language 和 Content-Type。
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		//设置是否允许携带凭据。这里设置为 "true",表示允许携带凭据。
		ctx.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		ctx.Next()
	}
}
