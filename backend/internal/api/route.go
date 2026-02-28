package api

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/moby/moby/api/types/events"
	"github.com/skye-z/harbor/internal/service"
	"github.com/skye-z/harbor/internal/util/docker"
	"github.com/skye-z/harbor/internal/util/response"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

// 路由结构体
type Route struct {
	Router       *gin.Engine
	DockerClient *docker.Client
	Engine       *xorm.Engine
}

// 创建路由实例
func NewRoute(page embed.FS) *Route {
	gin.SetMode(gin.ReleaseMode)
	route := &Route{}

	router, err := newRoute(page)
	if err != nil {
		log.Printf("[Core] 加载静态文件失败: %v", err)
		router = gin.Default()
	}
	route.Router = router

	client, err := docker.NewClient()
	if err != nil {
		log.Printf("[Core] 连接Docker失败: %v", err)
		route.DockerClient = nil
	} else {
		route.DockerClient = client
	}
	return route
}

// 创建引擎实例
func newRoute(page embed.FS) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	log.Println("[Core] 加载页面")
	pageFile, err := fs.Sub(page, "page/dist")
	if err != nil {
		log.Printf("[Core] 警告: 加载静态文件失败: %v", err)
		return router, nil
	}
	router.StaticFS("/app", http.FS(pageFile))
	return router, nil
}

// 初始化路由
func (r Route) Init(engine *xorm.Engine) {
	r.Engine = engine
	r.Router.Use(CORSMiddleware())

	if r.DockerClient != nil {
		ls := service.NewLogService(engine)
		service.SetLogEngine(engine)
		r.DockerClient.StartEventListener(func(event events.Message) {
			if event.Type == "container" {
				if event.Action == "start" || event.Action == "stop" || event.Action == "restart" || event.Action == "die" {
					containerName := ""
					if len(event.Actor.Attributes) > 0 {
						if name, ok := event.Actor.Attributes["name"]; ok {
							containerName = name
						}
					}
					if containerName == "" {
						containerName = event.Actor.ID
					}
					containerID := event.Actor.ID
					if len(containerID) > 12 {
						containerID = containerID[:12]
					}
					log.Printf("[Docker] docker container %s %s, name=%s", containerID, event.Action, containerName)
					err := ls.LogContainerEvent(string(event.Action), containerName, containerID, "")
					if err != nil {
						log.Printf("[Core] Failed to log container event: %v", err)
					} else {
						log.Printf("[Core] Container event logged successfully")
					}
				}
			} else if event.Type == "daemon" {
				if event.Action == "shutdown" {
					log.Println("[Docker] docker daemon shutdown")
				}
			}
		})
	}

	r.addPublicRoute()
	private := r.Router.Group("")
	private.Use(AuthMiddleware())
	{
		r.addPrivateRoute(private)
	}
}

// 添加公共路由
func (r Route) addPublicRoute() {
	r.Router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/app"
		r.Router.HandleContext(ctx)
	})

	as := service.NewAuthService(r.Engine)
	ls := service.NewLogService(r.Engine)
	r.Router.POST("/api/user/login", func(c *gin.Context) {
		var req service.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误")
			return
		}
		resp, err := as.Login(&req)
		if err != nil {
			response.Unauthorized(c, err.Error())
			return
		}
		ipAddress := c.ClientIP()
		ls.LogLogin(resp.Username, resp.UserID, ipAddress)
		response.Success(c, resp)
	})
}

// 添加私有路由
func (r Route) addPrivateRoute(route gin.IRoutes) {
	ds := service.NewDockerService(r.DockerClient)
	route.GET("/api/docker/info", ds.GetInfo)
	route.GET("/api/docker/close", ds.Close)
	route.GET("/api/docker/prune/containers", ds.PruneContainers)
	route.GET("/api/docker/prune/images", ds.PruneImages)
	route.GET("/api/docker/prune/volumes", ds.PruneVolumes)
	route.GET("/api/docker/prune/networks", ds.PruneNetworks)
	route.GET("/api/docker/prune/all", ds.PruneAll)

	ls := service.NewLogService(r.Engine)
	route.GET("/api/logs/recent", ls.GetRecent)

	is := service.NewImageService(r.DockerClient)
	route.GET("/api/image/list", is.GetList)
	route.GET("/api/image/search", is.SearchImages)
	route.GET("/api/image/pull", is.PullImage)
	route.GET("/api/image/pull/progress", is.GetPullProgress)
	route.GET("/api/image/remove", is.RemoveImage)
	route.GET("/api/image/inspect", is.GetInspect)
	route.GET("/api/image/build", is.BuildImage)
	route.GET("/api/image/tag", is.TagImage)
	route.GET("/api/image/push", is.PushImage)

	ns := service.NewNetworkService(r.DockerClient)
	route.GET("/api/network/list", ns.GetList)
	route.GET("/api/network/create", ns.Create)
	route.GET("/api/network/remove", ns.Remove)
	route.GET("/api/network/connect", ns.ConnectContainer)
	route.GET("/api/network/disconnect", ns.DisconnectContainer)

	vs := service.NewVolumeService(r.DockerClient)
	route.GET("/api/volume/list", vs.GetList)
	route.GET("/api/volume/create", vs.Create)
	route.GET("/api/volume/remove", vs.Remove)

	cs := service.NewContainerService(r.DockerClient)
	route.GET("/api/container/list", cs.GetList)
	route.GET("/api/container/info", cs.GetInfo)
	route.GET("/api/container/logs", cs.GetLogs)
	route.GET("/api/container/operation", cs.Operation)
	route.GET("/api/container/stat", cs.GetStat)
	route.GET("/api/container/processes", cs.GetProcesses)
	route.GET("/api/container/diff", cs.GetDiff)
	route.GET("/api/container/files", cs.ListFiles)
	route.GET("/api/container/copy/from", cs.CopyFromContainer)
	route.POST("/api/container/copy/to", cs.CopyToContainer)
	route.GET("/api/container/terminal", cs.ConnectTerminal)
	route.GET("/api/container/terminal/ws", cs.TerminalWebSocket)
	route.GET("/api/container/terminal/resize", cs.ResizeTerminal)
	route.GET("/api/container/terminal/close", cs.CloseTerminal)
	route.GET("/api/container/commit", cs.CommitContainer)
	route.POST("/api/container/create", cs.CreateContainer)
	route.POST("/api/container/rename", cs.RenameContainer)

	us := service.NewUserService(r.Engine)
	route.GET("/api/user/list", us.GetList)
	route.POST("/api/user/create", us.Create)
	route.POST("/api/user/update", us.Update)
	route.GET("/api/user/delete", us.Delete)
	route.GET("/api/user/current", us.GetCurrentUser)
}

// 获取服务端口
func (r Route) GetPort() string {
	if viper.IsSet("server.port") {
		return strconv.Itoa(viper.GetInt("server.port"))
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "12800"
	}
	return port
}

// 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Upgrade") == "websocket" {
			c.Next()
			return
		}
		token := c.GetHeader("Authorization")
		if token != "" && len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			response.Unauthorized(c, "未授权访问")
			c.Abort()
			return
		}

		payload, err := service.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, "无效的认证令牌")
			c.Abort()
			return
		}
		c.Set("user_id", payload.UserID)
		c.Set("username", payload.Username)
		c.Set("is_admin", payload.IsAdmin)
		c.Next()
	}
}

// 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigin := c.Request.Header.Get("Origin")
		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:12800")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With, Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// 压缩
func GzipMiddleware() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}

// 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

// 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	})
}

// 超时中间件
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("timeout_start", time.Now())
		c.Next()
	}
}

// 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// 生成请求ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// 生成随机字符串
func randomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)[:length]
}
