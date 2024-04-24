package router

import (
	"github.com/gin-gonic/gin"
)

var initRouterFunc = make([]func(*gin.RouterGroup), 0)

func Init(routerGroup *gin.RouterGroup) {
	for _, f := range initRouterFunc {
		f(routerGroup)
	}
	//initSysUser(routerGroup)
	//initSystem(routerGroup)
	//initSysDict(routerGroup)
	//initSysMenu(routerGroup)
	//initSysRole(routerGroup)
	//initSysPost(routerGroup)
	//initSysLoginPolicy(routerGroup)
	//initAssetNetwork(routerGroup)
	//initAssetInfo(routerGroup)
	//initAssetGroup(routerGroup)
}
