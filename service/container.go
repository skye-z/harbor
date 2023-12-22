package service

import (
	"harbor/docker"
	"harbor/util"

	"github.com/gin-gonic/gin"
)

type ContainerService struct {
	Client *docker.Docker
}

func NewContainerService(client *docker.Docker) *ContainerService {
	ds := new(ContainerService)
	ds.Client = client
	return ds
}

func (ds ContainerService) GetList(ctx *gin.Context) {
	list, err := ds.Client.GetContainerList()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

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
	timeout := 10000
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
	timeout := 10000
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
