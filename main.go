package main

import (
	"embed"
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
	route.Init(engine)

	port := route.GetPort()
	log.Println("[Core] service started, port is", port)
	// 启动服务
	go func() {
		if err := route.Router.Run(":" + port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	waitForInterrupt()
}

func loadDBEngine() *xorm.Engine {
	log.Println("[Data] load engine")
	engine, err := xorm.NewEngine("sqlite", "./local.store")
	if err != nil {
		panic(err)
	}
	return engine
}

func waitForInterrupt() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("[Core] Shutting down server...")

	// 在此处执行清理和关闭操作
	// 例如，你可能需要关闭数据库连接等

	log.Println("[Core] Server gracefully stopped")
}
