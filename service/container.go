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
