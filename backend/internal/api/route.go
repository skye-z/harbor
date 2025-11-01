package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/skye-z/harbor/internal/util/docker"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type Route struct {
	Router       *gin.Engine
	DockerClient *docker.Client
}

// 创建路由
func NewRoute(page embed.FS) *Route {
	gin.SetMode(gin.ReleaseMode)
	route := new(Route)
	route.Router = newRoute(page)

	client, err := docker.NewClient()
	if err != nil {
		return nil
	}
	route.DockerClient = client
	return route
}

// 创建路由
func newRoute(page embed.FS) *gin.Engine {
	router := gin.Default()
	log.Println("[Core] load page")
	pageFile, _ := fs.Sub(page, "page/dist")
	router.StaticFS("/app", http.FS(pageFile))
	return router
}

// 初始化路由
func (r Route) Init(engine *xorm.Engine) {
	// us := service.NewUserService(engine)
	r.addPublicRoute()
	// 私有路由
	private := r.Router.Group("")
	// .Use(service.AuthHandler())
	{
		r.addPrivateRoute(private)
	}
}

// 公共路由
func (r Route) addPublicRoute() {
	r.Router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/app"
		r.Router.HandleContext(ctx)
	})

	// r.Router.POST("/api/user/login", us.Login)
	// r.Router.GET("/api/oauth2/state", us.State)
}

// 私有路由
func (r Route) addPrivateRoute(route gin.IRoutes) {
	// Docker 系统管理路由
	// route.GET("/api/docker/close", ds.Close)
	// route.GET("/api/docker/prune/containers", ds.PruneContainers)
	// route.GET("/api/docker/prune/images", ds.PruneImages)
	// route.GET("/api/docker/prune/volumes", ds.PruneVolumes)
	// route.GET("/api/docker/prune/networks", ds.PruneNetworks)
	// route.GET("/api/docker/prune/all", ds.PruneAll)
}

// 获取端口号配置
func (r Route) GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "12800"
	}
	return port
}
