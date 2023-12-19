package service

import (
	"harbor/docker"
	"harbor/util"

	"github.com/gin-gonic/gin"
)

type VolumeService struct {
	Client *docker.Docker
}

func NewVolumeService(client *docker.Docker) *VolumeService {
	ds := new(VolumeService)
	ds.Client = client
	return ds
}

func (ds VolumeService) GetList(ctx *gin.Context) {
	list, err := ds.Client.GetVolumeList()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}
