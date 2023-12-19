package service

import (
	"harbor/docker"
	"harbor/util"

	"github.com/gin-gonic/gin"
)

type ImageService struct {
	Client *docker.Docker
}

func NewImageService(client *docker.Docker) *ImageService {
	ds := new(ImageService)
	ds.Client = client
	return ds
}

func (ds ImageService) GetList(ctx *gin.Context) {
	list, err := ds.Client.GetImageList()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}
