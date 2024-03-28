package middleware

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
	"go-protector/server/core/consts"
	"go-protector/server/core/current"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_jwt"
	"go-protector/server/core/custom/c_logger"
	"net/http"
	"path"
	"strings"
)

var ignoreUrlSet map[string]map[string]any

func init() {
	ignoreUrlSet = map[string]map[string]any{
		"POST": {
			path.Join(consts.ServerUrlPrefix, "/user", "login"): consts.EmptyVal,
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
				c_logger.GetLogger(c).Debug("Hit whitelist break auth, method: %s, url: %s", c.Request.Method, c.Request.URL.Path)
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
		if split := strings.Split(tokenStr, " "); len(split) >= 1 {
			tokenStr = split[1]
		}
		oldTokenStr := tokenStr
		currentUser, err := c_jwt.ParserToken(&tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, base.ResultCustom(http.StatusUnauthorized, nil, err.Error()))
			return
		}
		c.Request = c.Request.WithContext(current.SetUser(c.Request.Context(), currentUser))
		if oldTokenStr != tokenStr {
			c.Writer.Header().Set(consts.AuthHeaderKey, tokenStr)
			c_logger.GetLogger(c).Info("-------------------")
			c_logger.GetLogger(c).Info("old: %s", oldTokenStr)
			c_logger.GetLogger(c).Info("new: %s", tokenStr)
			c_logger.GetLogger(c).Info("-------------------")
		}

	}
}
