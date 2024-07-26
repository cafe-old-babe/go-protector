package middleware

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_jwt"
	"go-protector/server/internal/custom/c_logger"
	"net/http"
	"path"
	"strings"
)

// 3-10	【实战】JWT主动销毁及续约（忽略白名单校验，并发场景下续约的两种解决方案分析与实现）
var ignoreUrlSet map[string]map[string]any

func init() {
	ignoreUrlSet = map[string]map[string]any{
		"POST": {
			path.Join(consts.ServerUrlPrefix, "/user", "login"):  consts.EmptyVal,
			path.Join(consts.ServerUrlPrefix, "/user", "logout"): consts.EmptyVal,
		},
		"GET": {
			path.Join(consts.ServerUrlPrefix, "system", "/captcha"): consts.EmptyVal,
		},
	}
}

// JwtAuth 身份认证，包括换 token
func JwtAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		if set, ok := ignoreUrlSet[c.Request.Method]; ok {
			if _, ok := set[c.FullPath()]; ok {
				c_logger.GetLoggerByCtx(c).Debug("Hit whitelist break auth, method: %s, url: %s", c.Request.Method, c.Request.URL.Path)
				c.Next()
				return
			}
		}
		var tokenStr string
		if tokenStr = c.Request.Header.Get(consts.AuthHeaderKey); len(tokenStr) <= 0 {
			if tokenStr = c.Query(consts.AuthUrlKey); len(tokenStr) <= 0 {
				_ = c.AbortWithError(http.StatusUnauthorized, c_error.ErrAuthFailure)
				return
			}
		}
		if split := strings.Split(tokenStr, " "); len(split) >= 2 {
			tokenStr = split[1]
		}
		oldTokenStr := tokenStr
		currentUser, err := c_jwt.ParserToken(&tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, base.ResultCustom(http.StatusUnauthorized, nil, err.Error()))
			return
		}
		// 4-4	【实战】实现ant-design-vue-pro的路由接口-掌握使用Gin中间件保存当前用户信息
		c.Request = c.Request.WithContext(current.SetUser(c.Request.Context(), currentUser))

		if oldTokenStr != tokenStr {
			//c.Header(consts.AuthHeaderKey, tokenStr)
			c.Writer.Header().Set(consts.AuthHeaderKey, tokenStr)

		}
		c.Next()
	}
}
