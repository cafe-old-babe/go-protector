package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/core/config"
	"go-protector/server/core/consts"
	"go-protector/server/router"
	"go-protector/server/router/middleware"
	"net/http"
)

func initServer() (server *http.Server) {
	engine := gin.New()
	gin.SetMode(config.GetConfig().Server.Model)
	engine.Use(middleware.Cors())
	engine.Use(middleware.Recovery)
	engine.Use(middleware.JwtAuth())
	engine.Use(middleware.RecordLog)
	engine.Use(middleware.SetDB)

	routerGroup := engine.Group(consts.ServerUrlPrefix)
	router.Init(routerGroup)

	//if err = engine.Run(fmt.Sprintf(":%d", config._config.Server.Port)); err != nil {
	//	return err
	//}
	server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.GetConfig().Server.Port),
		Handler: engine,
	}
	return
}
