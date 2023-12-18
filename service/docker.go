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
