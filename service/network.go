package service

import (
	"harbor/docker"
	"harbor/util"

	"github.com/gin-gonic/gin"
)

type NetworkService struct {
	Client *docker.Docker
}

func NewNetworkService(client *docker.Docker) *NetworkService {
	ds := new(NetworkService)
	ds.Client = client
	return ds
}

func (ds NetworkService) GetList(ctx *gin.Context) {
	list, err := ds.Client.GetNetworkList()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}
