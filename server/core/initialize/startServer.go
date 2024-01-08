package initialize

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// StartServer https://gin-gonic.com/zh-cn/docs/examples/
func StartServer() (err error) {
	// 加载配置
	if err = initLogger(); err != nil {
		return
	}

	// 初始化并赋值数据库全局变量
	if err = initDB(); err != nil {
		return
	}
	// 初始化redis
	if err = initCache(); err != nil {
		return
	}

	server := initServer()
	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Error("ListenAndServe failure ", zap.Error(err))
			return
		}
	}()

	zap.L().Info("server start success!")
	// 优雅停机
	// 停止进程信号主要是3个，SIGINT,SIGKILL,SIGTERM。
	// SIGINT
	//产生方式为键盘ctrl+c，只针对当前前台进程及其所在的进程组的每个进程都发送SIGINT信号，之后这些进程会执行信号处理程序再终止。
	//
	//SIGTERM
	//产生方式为使用kill函数，kill + pid 方式发送。当前进程会收到信号，而子进程不会收到，如果当前进程被kill，则其子进程的父进程将为init，即pid为1的进程。
	//
	//SIGKILL
	//产生方式为使用kill函数，kill -9 + pid方式发送，当前进程收到该信号，会直接发送默认行为（该信号无法捕获，即无法处理）。

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zap.L().Error("server Shutdown failure ", zap.Error(err))
	}

	return
}
