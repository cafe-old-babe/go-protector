package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/config"
	"go-protector/server/router"
	"go-protector/server/router/middleware"
	"net/http"
)

func initServer() (server *http.Server) {
	engine := gin.New()
	gin.SetMode(config.GetConfig().Server.Model)
	engine.Use(middleware.RecordLog)
	engine.Use(middleware.Recovery)
	engine.Use(middleware.SetDB)
	engine.Use(middleware.Cors())

	routerGroup := engine.Group("api")
	router.InitLogin(routerGroup)

	//if err = engine.Run(fmt.Sprintf(":%d", config._config.Server.Port)); err != nil {
	//	return err
	//}
	server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.GetConfig().Server.Port),
		Handler: engine,
	}
	return
}
