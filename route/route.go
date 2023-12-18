package route

import (
	"harbor/docker"
	"harbor/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Router       *gin.Engine
	DockerClient *docker.Docker
}

// 创建外部路由
func NewRoute() *Route {
	// 关闭调试
	gin.SetMode(gin.ReleaseMode)
	route := new(Route)
	route.Router = newRoute()

	client, err := docker.NewDocker()
	if err != nil {
		return nil
	}
	route.DockerClient = client
	return route
}

// 创建路由
func newRoute() *gin.Engine {
	router := gin.Default()
	// 加载错误模板
	// templ := template.Must(template.New("").ParseFS(page, "page/error/*.html"))
	// router.SetHTMLTemplate(templ)
	// 配置404错误页面
	// router.NoRoute(func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "404.html", gin.H{
	// 		"title": "404",
	// 	})
	// })
	// 加载页面
	// log.Println("[Core] Load page")
	// pageFile, _ := fs.Sub(page, "page/dist")
	// router.StaticFS("/app", http.FS(pageFile))
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
	dc := service.NewDockerService(r.DockerClient)
	// 获取系统配置
	r.Router.GET("/api/docker/info", dc.GetInfo)
}

// 私有路由
func (r Route) addPrivateRoute() {
}

// 启动路由
func (r Route) Run() {
	port := r.getPort()
	log.Println("[Core] service started, port is", port)
	r.Router.Run(":" + port)
}

// 获取端口号配置
func (r Route) getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1280"
	}
	return port
}
