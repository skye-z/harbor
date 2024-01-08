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
	vs := new(VolumeService)
	vs.Client = client
	return vs
}

func (vs VolumeService) GetList(ctx *gin.Context) {
	list, err := vs.Client.GetVolumeList()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

func (vs VolumeService) GetInfo(ctx *gin.Context) {
	id := ctx.Query("id")
	list, err := vs.Client.GetVolumeInfo(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取存储详情失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

func (vs VolumeService) Create(ctx *gin.Context) {
	name := ctx.Query("name")
	err := vs.Client.CreateVolume(name)
	if err != nil {
		util.ReturnMessage(ctx, false, "创建存储失败")
	} else {
		util.ReturnData(ctx, true, "创建存储成功")
	}
}

func (vs VolumeService) Remove(ctx *gin.Context) {
	id := ctx.Query("id")
	err := vs.Client.RemoveVolume(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "删除存储失败")
	} else {
		util.ReturnMessage(ctx, true, "删除存储成功")
	}
}
