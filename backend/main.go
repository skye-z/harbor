package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/skye-z/harbor/internal/data"
	"github.com/skye-z/harbor/internal/util/config"
	"xorm.io/xorm"
)

func main() {
	// 初始化系统配置
	config.InitConfig()

	// 初始化数据库
	// 初始化数据库
	engine := data.InitDB()
	go data.InitDBTable(engine)

	// 设置 GIN
	// TODO 下面是从v1 copy的
	// route := route.NewRoute(page)
	// if route == nil {
	// 	log.Println("[Core] please start docker first")
	// 	return
	// }
	// route.Init(engine)
	// // 获取端口
	// port := route.GetPort()
	// log.Println("[Core] service started, port is", port)
	// // 启动服务
	// go func() {
	// 	if err := route.Router.Run(":" + port); err != nil {
	// 		log.Fatalf("Error starting server: %v", err)
	// 	}
	// }()

	// 等待中断信号以优雅关闭服务器
	waitForInterrupt(engine)
}

func waitForInterrupt(engine *xorm.Engine) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("[Core] Shutting down server...")

	defer engine.Close()

	log.Println("[Core] Server gracefully stopped")
}
