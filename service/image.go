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
	is := new(ImageService)
	is.Client = client
	return is
}

// 获取镜像列表
func (is ImageService) GetList(ctx *gin.Context) {
	list, err := is.Client.GetImageList()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取镜像列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

// 删除镜像
func (is ImageService) Remove(ctx *gin.Context) {
	id := ctx.Query("id")
	prune := ctx.Query("prune")
	force := ctx.Query("id")
	err := is.Client.RemoveImage(id, prune == "1", force == "1")
	if err != nil {
		util.ReturnMessage(ctx, false, "镜像删除失败")
	} else {
		util.ReturnMessage(ctx, true, "镜像删除成功")
	}
}

// 拉取镜像
func (is ImageService) Pull(ctx *gin.Context) {
	var form docker.ImageBuild
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if len(form.Store) == 0 {
		util.ReturnMessage(ctx, false, "镜像仓库不能为空")
		return
	}
	if len(form.Tag) == 0 {
		util.ReturnMessage(ctx, false, "镜像标签不能为空")
		return
	}
	if len(form.Platform) == 0 {
		util.ReturnMessage(ctx, false, "平台架构不能为空")
		return
	}
	is.Client.PullImage(ctx, form)
}

// 打标签
func (is ImageService) AddTag(ctx *gin.Context) {
	id := ctx.Query("id")
	tag := ctx.Query("tag")
	err := is.Client.AddImageTag(id, tag)
	if err != nil {
		util.ReturnMessage(ctx, false, "镜像标记失败")
	} else {
		util.ReturnMessage(ctx, true, "镜像标记成功")
	}
}

// 获取镜像详情
func (is ImageService) GetInfo(ctx *gin.Context) {
	id := ctx.Query("id")
	info, err := is.Client.GetImageInfo(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取镜像详情失败")
	} else {
		util.ReturnData(ctx, true, info)
	}
}

// 获取镜像构建历史
func (is ImageService) GetHistory(ctx *gin.Context) {
	id := ctx.Query("id")
	info, err := is.Client.GetImageHistory(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取镜像历史失败")
	} else {
		util.ReturnData(ctx, true, info)
	}
}
