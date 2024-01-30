package service

import (
	"harbor/docker"
	"harbor/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 连接升级程序
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ContainerService struct {
	Client *docker.Docker
}

func NewContainerService(client *docker.Docker) *ContainerService {
	ds := new(ContainerService)
	ds.Client = client
	return ds
}

// 获取容器列表
func (ds ContainerService) GetList(ctx *gin.Context) {
	list, err := ds.Client.GetContainerList()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

// 获取容器详情
func (ds ContainerService) GetInfo(ctx *gin.Context) {
	id := ctx.Query("id")
	info, err := ds.Client.GetContainerInfo(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取详情失败")
	} else {
		util.ReturnData(ctx, true, info)
	}
}

// 启动容器
func (ds ContainerService) StartContainer(ctx *gin.Context) {
	id := ctx.Query("id")
	err := ds.Client.StartContainer(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "容器启动失败")
	} else {
		util.ReturnMessage(ctx, true, "容器启动成功")
	}
}

// 重启容器
func (ds ContainerService) RestartContainer(ctx *gin.Context) {
	id := ctx.Query("id")
	timeout := 30
	err := ds.Client.RestartContainer(id, &timeout)
	if err != nil {
		util.ReturnMessage(ctx, false, "容器重启失败")
	} else {
		util.ReturnMessage(ctx, true, "容器重启成功")
	}
}

// 停止容器
func (ds ContainerService) StopContainer(ctx *gin.Context) {
	id := ctx.Query("id")
	timeout := 30
	err := ds.Client.StopContainer(id, &timeout)
	if err != nil {
		util.ReturnMessage(ctx, false, "容器停止失败")
	} else {
		util.ReturnMessage(ctx, true, "容器停止成功")
	}
}

// 关闭容器
func (ds ContainerService) KillContainer(ctx *gin.Context) {
	id := ctx.Query("id")
	err := ds.Client.KillContainer(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "容器关闭失败")
	} else {
		util.ReturnMessage(ctx, true, "容器关闭成功")
	}
}

// 挂起容器
func (ds ContainerService) PauseContainer(ctx *gin.Context) {
	id := ctx.Query("id")
	err := ds.Client.PauseContainer(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "容器挂起失败")
	} else {
		util.ReturnMessage(ctx, true, "容器挂起成功")
	}
}

// 从挂起中恢复容器
func (ds ContainerService) UnpauseContainer(ctx *gin.Context) {
	id := ctx.Query("id")
	err := ds.Client.UnpauseContainer(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "容器恢复失败")
	} else {
		util.ReturnMessage(ctx, true, "容器恢复成功")
	}
}

// 删除容器
func (ds ContainerService) RemoveContainer(ctx *gin.Context) {
	id := ctx.Query("id")
	volumes := ctx.Query("volumes")
	links := ctx.Query("links")
	force := ctx.Query("force")
	err := ds.Client.RemoveContainer(id, volumes == "1", links == "1", force == "1")
	if err != nil {
		util.ReturnMessage(ctx, false, "容器删除失败")
	} else {
		util.ReturnMessage(ctx, true, "容器删除成功")
	}
}

// 获取容器日志
func (ds ContainerService) GetLogs(ctx *gin.Context) {
	id := ctx.Query("id")
	tail := ctx.DefaultQuery("tail", "100")
	logs, err := ds.Client.GetContainerLogs(id, tail)
	if err != nil {
		util.ReturnMessage(ctx, false, "容器日志读取失败")
	} else {
		util.ReturnData(ctx, true, logs)
	}
}

// 连接容器终端
func (ds ContainerService) ConnectTerminal(ctx *gin.Context) {
	id := ctx.Query("id")
	cmd := ctx.DefaultQuery("cmd", "/bin/sh")
	cols, _ := strconv.Atoi(ctx.DefaultQuery("cols", "80"))
	rows, _ := strconv.Atoi(ctx.DefaultQuery("rows", "90"))

	// 升级连接
	upgrade, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	defer upgrade.Close()

	err = ds.Client.CreateTerminal(upgrade, id, cmd, uint(cols), uint(rows))
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
}

// 获取容器变动
func (ds ContainerService) GetDiff(ctx *gin.Context) {
	id := ctx.Query("id")
	list, err := ds.Client.GetContainerDiff(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取容器变动失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

// 获取容器统计信息
func (ds ContainerService) GetStat(ctx *gin.Context) {
	id := ctx.Query("id")
	stat, err := ds.Client.GetContainerStat(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取容器统计信息失败")
	} else {
		util.ReturnData(ctx, true, stat)
	}
}

// 获取容器进程信息
func (ds ContainerService) GetProcesses(ctx *gin.Context) {
	id := ctx.Query("id")
	list, err := ds.Client.GetContainerProcesses(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取容器进程信息失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

// 克隆容器
func (ds ContainerService) CloneContainer(ctx *gin.Context) {
	id := ctx.Query("id")
	res, err := ds.Client.CloneContainer(id)
	if err != nil {
		util.ReturnMessageData(ctx, false, "容器克隆失败", err.Error())
	} else {
		util.ReturnMessageData(ctx, true, "容器克隆成功", res)
	}
}

// 重建容器
func (ds ContainerService) RecreateContainer(ctx *gin.Context) {
	id := ctx.Query("id")
	res, err := ds.Client.RecreateContainer(id)
	if err != nil {
		util.ReturnMessageData(ctx, false, "容器重建失败", err.Error())
	} else {
		util.ReturnMessageData(ctx, true, "容器重建成功", res)
	}
}
