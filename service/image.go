/*
镜像服务

BetaX Harbor
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package service

import (
	"log"
	"strconv"

	"github.com/skye-z/harbor/docker"
	"github.com/skye-z/harbor/model"
	"github.com/skye-z/harbor/util"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

const (
	DATA_ERROR = "传入数据无效"
	TEMP_ERROR = "请先指定模板"
)

type ImageService struct {
	Client     *docker.Docker
	ImageModel model.ImageModel
}

func NewImageService(client *docker.Docker, engine *xorm.Engine) *ImageService {
	is := new(ImageService)
	is.Client = client
	is.ImageModel = model.ImageModel{
		DB: engine,
	}
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
		util.ReturnMessage(ctx, false, DATA_ERROR)
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

// 搜索镜像
func (is ImageService) SearchImage(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	list, err := is.Client.SearchImage(keyword)
	if err != nil {
		util.ReturnMessage(ctx, false, "搜索镜像失败")
	} else {
		util.ReturnData(ctx, true, list)
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

// 导出镜像
func (is ImageService) ExportImage(ctx *gin.Context) {
	id := ctx.Query("id")
	is.Client.ExportImage(ctx, id)
}

// 构建镜像
func (is ImageService) BuildImage(ctx *gin.Context) {
	var form docker.ImageBuild
	if err := ctx.ShouldBindJSON(&form); err != nil {
		util.ReturnMessage(ctx, false, DATA_ERROR)
		return
	}
	if len(form.Tag) == 0 {
		util.ReturnMessage(ctx, false, "镜像标签不能为空")
		return
	}
	if len(form.Context) == 0 {
		util.ReturnMessage(ctx, false, "镜像构建内容不能为空")
		return
	}
	log, err := is.Client.BuildImage(form.Tag, form.Context)
	if err != nil {
		util.ReturnMessageData(ctx, false, "镜像构建失败", err.Error())
	} else {
		util.ReturnMessageData(ctx, true, "镜像构建成功", log)
	}
}

// 添加构建模板
func (is ImageService) AddImageBuild(ctx *gin.Context) {
	var form model.Image
	if err := ctx.ShouldBindJSON(&form); err != nil {
		log.Println(err)
		util.ReturnMessage(ctx, false, DATA_ERROR)
		return
	}
	if len(form.Name) == 0 {
		util.ReturnMessage(ctx, false, "构建模板名称不能为空")
		return
	}
	if len(form.Context) == 0 {
		util.ReturnMessage(ctx, false, "构建模板内容不能为空")
		return
	}
	state := is.ImageModel.AddImage(form.Name, form.Context)
	if state {
		util.ReturnMessage(ctx, true, "模板添加成功")
	} else {
		util.ReturnMessage(ctx, false, "模板添加失败")
	}
}

// 编辑构建模板
func (is ImageService) EditImageBuild(ctx *gin.Context) {
	var form model.Image
	if err := ctx.ShouldBindJSON(&form); err != nil {
		util.ReturnMessage(ctx, false, DATA_ERROR)
		return
	}
	if len(form.Name) == 0 {
		util.ReturnMessage(ctx, false, "构建模板名称不能为空")
		return
	}
	if len(form.Context) == 0 {
		util.ReturnMessage(ctx, false, "构建模板内容不能为空")
		return
	}

	state := is.ImageModel.EditImage(&form)
	if state {
		util.ReturnMessage(ctx, true, "模板编辑成功")
	} else {
		util.ReturnMessage(ctx, false, "模板编辑失败")
	}
}

// 删除构建模板
func (is ImageService) DelImageBuild(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		util.ReturnMessage(ctx, false, TEMP_ERROR)
		return
	}
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.ReturnMessage(ctx, false, TEMP_ERROR)
		return
	}

	state := is.ImageModel.DelImage(tid)
	if state {
		util.ReturnMessage(ctx, true, "模板删除成功")
	} else {
		util.ReturnMessage(ctx, false, "模板删除失败")
	}
}

// 获取构建模板列表
func (is ImageService) GetImageBuildList(ctx *gin.Context) {
	list, err := is.ImageModel.GetImageList()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取构建模板列表失败")
		return
	}
	util.ReturnData(ctx, true, list)
}

// 获取构建模板详情
func (is ImageService) GetImageBuildInfo(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		util.ReturnMessage(ctx, false, TEMP_ERROR)
		return
	}
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.ReturnMessage(ctx, false, TEMP_ERROR)
		return
	}

	info, err := is.ImageModel.GetImageInfo(tid)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取构建模板详情失败")
	} else {
		util.ReturnData(ctx, true, info)
	}
}
