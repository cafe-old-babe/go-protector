package router

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	initLogin(routerGroup)
	initSystem(routerGroup)
	initSysDict(routerGroup)
}
