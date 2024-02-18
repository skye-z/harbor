package main

import (
	"embed"
	"harbor/model"
	"harbor/monitor"
	"harbor/route"
	"harbor/service"
	"harbor/util"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

//go:embed page/dist/*
var page embed.FS

func main() {
	// 初始化系统配置
	util.InitConfig()
	// 初始化数据库
	engine := loadDBEngine()
	go service.InitDatabase(engine)
	route := route.NewRoute(page)
	if route == nil {
		log.Println("[Core] please start docker first")
		return
	}
	route.Init(engine)
	// 启动监控
	runMonitor(engine)
	// 获取端口
	port := route.GetPort()
	log.Println("[Core] service started, port is", port)
	// 写入启动日志
	addMainLog(engine, "start")
	// 启动服务
	go func() {
		if err := route.Router.Run(":" + port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	waitForInterrupt(engine)
}

func loadDBEngine() *xorm.Engine {
	log.Println("[Data] load engine")
	engine, err := xorm.NewEngine("sqlite", "./local.store")
	if err != nil {
		panic(err)
	}
	return engine
}

func runMonitor(engine *xorm.Engine) {
	go monitor.ListenDockerEvents(engine)
	go monitor.ListenHostOverhead(engine)
}

func waitForInterrupt(engine *xorm.Engine) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("[Core] Shutting down server...")

	addMainLog(engine, "stop")
	defer engine.Close()

	log.Println("[Core] Server gracefully stopped")
}

func addMainLog(engine *xorm.Engine, action string) {
	logger := &model.LogModel{
		DB: engine,
	}
	logger.AddLog("platform", action, "")
}
