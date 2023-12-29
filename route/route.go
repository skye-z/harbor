package route

import (
	"embed"
	"harbor/docker"
	"harbor/service"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Router       *gin.Engine
	DockerClient *docker.Docker
}

// 创建外部路由
func NewRoute(page embed.FS) *Route {
	// 关闭调试
	gin.SetMode(gin.ReleaseMode)
	route := new(Route)
	route.Router = newRoute(page)

	client, err := docker.NewDocker()
	if err != nil {
		return nil
	}
	route.DockerClient = client
	return route
}

// 创建路由
func newRoute(page embed.FS) *gin.Engine {
	router := gin.Default()
	// 加载页面
	log.Println("[Core] load page")
	pageFile, _ := fs.Sub(page, "page/dist")
	router.StaticFS("/app", http.FS(pageFile))
	return router
}

// 初始化路由
func (r Route) Init() {
	r.addPublicRoute()
	r.addPrivateRoute()
}

// 公共路由
func (r Route) addPublicRoute() {
	r.Router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/app"
		r.Router.HandleContext(ctx)
	})
	ds := service.NewDockerService(r.DockerClient)
	r.Router.GET("/api/docker/info", ds.GetInfo)

	cs := service.NewContainerService(r.DockerClient)
	r.Router.GET("/api/container/list", cs.GetList)
	r.Router.GET("/api/container/info", cs.GetInfo)
	r.Router.GET("/api/container/logs", cs.GetLogs)
	r.Router.GET("/api/container/start", cs.StartContainer)
	r.Router.GET("/api/container/stop", cs.StopContainer)
	r.Router.GET("/api/container/restart", cs.RestartContainer)
	r.Router.GET("/api/container/kill", cs.KillContainer)
	r.Router.GET("/api/container/pause", cs.PauseContainer)
	r.Router.GET("/api/container/unpause", cs.UnpauseContainer)
	r.Router.GET("/api/container/terminal", cs.ConnectTerminal)
	r.Router.GET("/api/container/diff", cs.GetDiff)
	r.Router.GET("/api/container/stat", cs.GetStat)
	r.Router.GET("/api/container/processes", cs.GetProcesses)

	is := service.NewImageService(r.DockerClient)
	r.Router.GET("/api/image/list", is.GetList)

	ns := service.NewNetworkService(r.DockerClient)
	r.Router.GET("/api/network/list", ns.GetList)

	vs := service.NewVolumeService(r.DockerClient)
	r.Router.GET("/api/volume/list", vs.GetList)
}

// 私有路由
func (r Route) addPrivateRoute() {
}

// 获取端口号配置
func (r Route) GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "12800"
	}
	return port
}
