package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/api"
)

func InitLogin(group *gin.RouterGroup) {
	user := group.Group("user")
	{
		userApi := api.UserApi
		user.POST("login", userApi.Login)
		user.POST("logout", userApi.Logout)
		user.POST("setStatus", userApi.SetStatus)
	}
}
