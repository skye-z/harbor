package route

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/service"
	"github.com/skye-z/harbor/internal/util/docker"
	"xorm.io/xorm"
)

type Route struct {
	Router       *gin.Engine
	DockerClient *docker.Client
	Engine       *xorm.Engine
}

// 创建外部路由
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
	r.Engine = engine
	r.addPublicRoute()
	// 私有路由，使用认证中间件
	private := r.Router.Group("")
	private.Use(AuthMiddleware())
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

	// 用户登录
	as := service.NewAuthService(r.Engine)
	r.Router.POST("/api/user/login", func(c *gin.Context) {
		var req service.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "请求参数错误"})
			return
		}
		resp, err := as.Login(&req)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resp)
	})
}

// 私有路由
func (r Route) addPrivateRoute(route gin.IRoutes) {
	ds := service.NewDockerService(r.DockerClient)
	route.GET("/api/docker/info", ds.GetInfo)
	route.GET("/api/docker/close", ds.Close)
	route.GET("/api/docker/prune/containers", ds.PruneContainers)
	route.GET("/api/docker/prune/images", ds.PruneImages)
	route.GET("/api/docker/prune/volumes", ds.PruneVolumes)
	route.GET("/api/docker/prune/networks", ds.PruneNetworks)
	route.GET("/api/docker/prune/all", ds.PruneAll)

	cs := service.NewContainerService(r.DockerClient)
	route.GET("/api/container/list", cs.GetList)
	route.GET("/api/container/info", cs.GetInfo)
	route.GET("/api/container/logs", cs.GetLogs)
	route.GET("/api/container/operation", cs.Operation)
	route.GET("/api/container/stat", cs.GetStat)
	route.GET("/api/container/processes", cs.GetProcesses)
	route.GET("/api/container/diff", cs.GetDiff)
	route.GET("/api/container/copy/from", cs.CopyFromContainer)
	route.GET("/api/container/copy/to", cs.CopyToContainer)
	route.GET("/api/container/terminal", cs.ConnectTerminal)
	route.GET("/api/container/terminal/ws", cs.TerminalWebSocket)
	route.GET("/api/container/terminal/resize", cs.ResizeTerminal)
	route.GET("/api/container/terminal/close", cs.CloseTerminal)
	route.GET("/api/container/commit", cs.CommitContainer)
}

// 获取端口号配置
func (r Route) GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "12800"
	}
	return port
}
