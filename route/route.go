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
	"xorm.io/xorm"
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
func (r Route) Init(engine *xorm.Engine) {
	us := service.NewUserService(engine)
	r.addOAuth2Route(engine)
	r.addPublicRoute(engine, us)
	// 私有路由
	private := r.Router.Group("").Use(service.AuthHandler())
	{
		r.addPrivateRoute(private, engine, us)
	}
}

// 公共路由
func (r Route) addPublicRoute(engine *xorm.Engine, us *service.UserService) {
	r.Router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/app"
		r.Router.HandleContext(ctx)
	})

	r.Router.POST("/api/user/login", us.Login)
	r.Router.GET("/api/oauth2/state", us.State)
}

// 私有路由
func (r Route) addPrivateRoute(route gin.IRoutes, engine *xorm.Engine, us *service.UserService) {
	ms := service.NewMonintorService(r.DockerClient)
	route.GET("/api/device/info", ms.GetDeviceInfo)
	route.GET("/api/system/use", ms.GetUse)

	ds := service.NewDockerService(r.DockerClient)
	route.GET("/api/docker/info", ds.GetInfo)

	cs := service.NewContainerService(r.DockerClient)
	route.GET("/api/container/list", cs.GetList)
	route.GET("/api/container/info", cs.GetInfo)
	route.GET("/api/container/logs", cs.GetLogs)
	route.GET("/api/container/start", cs.StartContainer)
	route.GET("/api/container/stop", cs.StopContainer)
	route.GET("/api/container/restart", cs.RestartContainer)
	route.GET("/api/container/kill", cs.KillContainer)
	route.GET("/api/container/pause", cs.PauseContainer)
	route.GET("/api/container/unpause", cs.UnpauseContainer)
	route.GET("/api/container/remove", cs.RemoveContainer)
	route.GET("/api/container/terminal", cs.ConnectTerminal)
	route.GET("/api/container/diff", cs.GetDiff)
	route.GET("/api/container/stat", cs.GetStat)
	route.GET("/api/container/processes", cs.GetProcesses)

	is := service.NewImageService(r.DockerClient, engine)
	route.GET("/api/image/list", is.GetList)
	route.GET("/api/image/remove", is.Remove)
	route.POST("/api/image/pull", is.Pull)
	route.GET("/api/image/tag", is.AddTag)
	route.GET("/api/image/search", is.SearchImage)
	route.GET("/api/image/info", is.GetInfo)
	route.GET("/api/image/history", is.GetHistory)
	route.GET("/api/image/export", is.ExportImage)

	route.POST("/api/image/build/run", is.BuildImage)
	route.POST("/api/image/build/add", is.AddImageBuild)
	route.POST("/api/image/build/edit", is.EditImageBuild)
	route.GET("/api/image/build/list", is.GetImageBuildList)
	route.GET("/api/image/build/info", is.GetImageBuildInfo)

	ns := service.NewNetworkService(r.DockerClient)
	route.GET("/api/network/list", ns.GetList)
	route.GET("/api/network/info", ns.GetInfo)
	route.POST("/api/network/create", ns.Create)
	route.GET("/api/network/remove", ns.Remove)
	route.GET("/api/network/connect", ns.Connect)
	route.GET("/api/network/disconnect", ns.Disconnect)

	vs := service.NewVolumeService(r.DockerClient)
	route.GET("/api/volume/list", vs.GetList)
	route.GET("/api/volume/info", vs.GetInfo)
	route.GET("/api/volume/create", vs.Create)
	route.GET("/api/volume/remove", vs.Remove)

	ss := service.NewSettingService()
	route.GET("/api/setting/oauth2", ss.GetOAuth2Setting)
	route.POST("/api/setting/oauth2", ss.UpdateOAuth2Setting)
	route.GET("/api/setting/alarm", ss.GetAlarmSetting)
	route.POST("/api/setting/alarm", ss.UpdateAlarmSetting)

	ls := service.NewLogsService(engine)
	route.GET("/api/log/number", ls.GetNumber)
	route.GET("/api/log/list", ls.GetLogs)

	route.GET("/api/user/list", us.GetList)
	route.GET("/api/user/info", us.GetInfo)
	route.POST("/api/user/add", us.Add)
	route.POST("/api/user/edit", us.Edit)
	route.GET("/api/user/remove", us.Del)
	route.GET("/api/user/bind", us.Bind)
}

// 授权登陆路由
func (r Route) addOAuth2Route(engine *xorm.Engine) {
	as := service.NewAuthService(engine)
	if as != nil {
		r.Router.GET("/oauth2/bind", as.Bind)
		r.Router.GET("/oauth2/login", as.Login)
		r.Router.GET("/oauth2/callback", as.Callback)
	}
}

// 获取端口号配置
func (r Route) GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "12800"
	}
	return port
}
