package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/skye-z/harbor/internal/api"
	"github.com/skye-z/harbor/internal/data"
	"github.com/skye-z/harbor/internal/service"
	"github.com/skye-z/harbor/internal/util/config"
	"xorm.io/xorm"
)

//go:embed dist/*
var page embed.FS

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.InitConfigWithPath("/opt/harbor")

	engine, err := data.InitDB()
	if err != nil {
		log.Fatalf("[Core] 初始化数据库失败: %v", err)
	}

	err = data.InitDBTable(engine)
	if err != nil {
		log.Fatalf("[Core] 初始化数据库表失败: %v", err)
	}

	logService := service.NewLogService(engine)
	if err := logService.LogSystemStartup(); err != nil {
		log.Printf("[Core] 记录系统启动日志失败: %v", err)
	}

	route := api.NewRoute(page)
	if route == nil {
		logService.Log(service.LogTypeSystem, service.LogLevelError, "startup", "system", "", "Docker未运行，无法启动服务", "", 0)
		log.Println("[Core] 请先启动Docker")
		return
	}
	route.Init(engine)
	port := route.GetPort()
	log.Println("[Core] 服务已启动，端口为", port)
	logService.Log(service.LogTypeSystem, service.LogLevelInfo, "startup", "system", "", "Harbor服务已启动，端口: "+port, "", 0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		srv := &http.Server{Addr: ":" + port}
		go func() {
			if err := route.Router.Run(":" + port); err != nil && err != http.ErrServerClosed {
				log.Fatalf("[Core] 启动服务器失败: %v", err)
			}
		}()
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	waitForInterrupt(engine, logService, cancel)
}

func waitForInterrupt(engine *xorm.Engine, logService *service.LogService, cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("[Core] 正在关闭服务器...")
	cancel()

	if logService != nil {
		logService.LogSystemShutdown()
	}

	defer engine.Close()

	log.Println("[Core] 服务器已优雅关闭")
}
