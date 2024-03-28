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
	// https://github.com/gin-gonic/gin/pull/3367
	// https://qiita.com/hum_op/items/901093b8bc3078b8077b 小日子
	// https://github.com/gin-gonic/gin/blob/v1.9.1/context.go#L1229-L1232
	// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
	// https://pkg.go.dev/github.com/rumorshub/gin#section-readme
	// https://before80.github.io/go_docs/thirdPkg/gin/gin/
	engine.ContextWithFallback = true

	gin.SetMode(config.GetConfig().Server.Model)
	engine.Use(middleware.Cors())
	engine.Use(middleware.RecordLog)
	engine.Use(middleware.Recovery)
	engine.Use(middleware.JwtAuth())
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
