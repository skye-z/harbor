package service

import (
	"harbor/docker"
	"harbor/util"

	"github.com/gin-gonic/gin"
)

type DockerService struct {
	Client *docker.Docker
}

func NewDockerService(client *docker.Docker) *DockerService {
	ds := new(DockerService)
	ds.Client = client
	return ds
}

func (ds DockerService) GetInfo(ctx *gin.Context) {
	util.ReturnData(ctx, true, ds.Client.Config)
}

func (ds DockerService) GetUsage(ctx *gin.Context) {
	info, err := ds.Client.GetUsage()
	if err != nil {
		util.ReturnMessageData(ctx, false, "获取使用情况失败", err.Error())
	} else {
		util.ReturnData(ctx, true, info)
	}
}

func (ds DockerService) CleanBuildCache(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	err := ds.Client.CleanBuildCache()
	if err != nil {
		util.ReturnMessageData(ctx, false, "构建缓存清理失败", err.Error())
	} else {
		util.ReturnMessage(ctx, true, "构建缓存清理成功")
	}
}

func (ds DockerService) CleanImage(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	err := ds.Client.CleanImage()
	if err != nil {
		util.ReturnMessageData(ctx, false, "清理失败", err.Error())
	} else {
		util.ReturnMessage(ctx, true, "清理成功")
	}
}

func (ds DockerService) CleanNetworks(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	err := ds.Client.CleanNetworks()
	if err != nil {
		util.ReturnMessageData(ctx, false, "清理失败", err.Error())
	} else {
		util.ReturnMessage(ctx, true, "清理成功")
	}
}

func (ds DockerService) CleanVolumes(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	err := ds.Client.CleanVolumes()
	if err != nil {
		util.ReturnMessageData(ctx, false, "清理失败", err.Error())
	} else {
		util.ReturnMessage(ctx, true, "清理成功")
	}
}
