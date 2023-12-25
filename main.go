package main

import (
	"embed"
	"harbor/route"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//go:embed page/dist/*
var page embed.FS

func main() {
	route := route.NewRoute(page)
	route.Init()

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

func waitForInterrupt() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("[Core] Shutting down server...")

	// 在此处执行清理和关闭操作
	// 例如，你可能需要关闭数据库连接等

	log.Println("[Core] Server gracefully stopped")
}
