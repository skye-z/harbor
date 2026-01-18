package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/util/docker"
	"github.com/skye-z/harbor/internal/util/response"
)

// 镜像管理服务
type ImageService struct {
	client *docker.Client
}

// 创建镜像服务实例
func NewImageService(client *docker.Client) *ImageService {
	return &ImageService{client: client}
}

// 获取镜像列表
func (s *ImageService) GetList(c *gin.Context) {
	images, err := s.client.ListImages()
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, images)
}

// 搜索镜像
func (s *ImageService) SearchImages(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		query = c.Query("term")
	}
	if query == "" {
		response.BadRequest(c, "缺少搜索关键词")
		return
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	results, err := s.client.SearchImages(query, limit)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, results)
}

// 拉取镜像
func (s *ImageService) PullImage(c *gin.Context) {
	tag := c.Query("image")
	if tag == "" {
		tag = c.Query("tag")
	}
	if tag == "" {
		response.BadRequest(c, "缺少镜像标签")
		return
	}

	if err := s.client.PullImage(tag); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "镜像拉取成功", gin.H{"tag": tag})
}

// 删除镜像
func (s *ImageService) RemoveImage(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少镜像ID")
		return
	}

	if err := s.client.RemoveImage(id); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "镜像删除成功", gin.H{"id": id})
}

// 获取镜像详情
func (s *ImageService) GetInspect(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少镜像ID")
		return
	}

	result, err := s.client.InspectImage(id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, result)
}

// 构建镜像
func (s *ImageService) BuildImage(c *gin.Context) {
	imageName := c.Query("name")
	if imageName == "" {
		response.BadRequest(c, "缺少镜像名称")
		return
	}

	dockerfile := c.DefaultQuery("dockerfile", "FROM alpine\nRUN apk add --no-cache curl")
	if err := s.client.BuildImage(imageName, dockerfile); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "镜像构建成功", gin.H{"name": imageName})
}

// 为镜像添加标签
func (s *ImageService) TagImage(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少镜像ID")
		return
	}

	tag := c.Query("tag")
	if tag == "" {
		response.BadRequest(c, "缺少标签")
		return
	}

	if err := s.client.TagImage(id, tag); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "镜像标签添加成功", gin.H{
		"id":  id,
		"tag": tag,
	})
}

// 推送镜像
func (s *ImageService) PushImage(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		response.BadRequest(c, "缺少镜像标签")
		return
	}

	if err := s.client.PushImage(tag); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "镜像推送成功", gin.H{"tag": tag})
}
