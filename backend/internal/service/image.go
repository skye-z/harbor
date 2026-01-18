package service

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/util/docker"
	"github.com/skye-z/harbor/internal/util/response"
	"github.com/skye-z/harbor/internal/util/validation"
)

type ImageService struct {
	client *docker.Client
}

func NewImageService(client *docker.Client) *ImageService {
	return &ImageService{client: client}
}

var (
	currentPullTag      string
	currentPullProgress docker.PullProgress
	pullLayerProgress   = make(map[string]docker.PullProgress)
	pullMutex           sync.RWMutex
)

func (s *ImageService) GetList(c *gin.Context) {
	images, err := s.client.ListImages()
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, images)
}

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

type LayerProgress struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
	Current  int64  `json:"current"`
	Total    int64  `json:"total"`
}

type PullProgressResponse struct {
	Tag     string          `json:"tag"`
	Layers  []LayerProgress `json:"layers"`
	Percent int             `json:"percent"`
}

func (s *ImageService) PullImage(c *gin.Context) {
	tag := c.Query("image")
	if tag == "" {
		tag = c.Query("tag")
	}
	if tag == "" {
		response.BadRequest(c, "缺少镜像标签")
		return
	}

	pullMutex.Lock()
	currentPullTag = tag
	currentPullProgress = docker.PullProgress{Status: "pending"}
	pullLayerProgress = make(map[string]docker.PullProgress)
	pullMutex.Unlock()

	go func() {
		onProgress := func(progress docker.PullProgress) {
			pullMutex.Lock()
			if progress.ID != "" {
				pullLayerProgress[progress.ID] = progress
			}
			currentPullProgress = progress
			pullMutex.Unlock()
		}

		err := s.client.PullImage(tag, onProgress)

		pullMutex.Lock()
		if err != nil {
			currentPullProgress = docker.PullProgress{Status: "error"}
		} else {
			currentPullTag = ""
			currentPullProgress = docker.PullProgress{}
			pullLayerProgress = make(map[string]docker.PullProgress)
		}
		pullMutex.Unlock()
	}()

	response.SuccessWithMessage(c, "镜像拉取任务已启动", gin.H{"tag": tag})
}

func (s *ImageService) GetPullProgress(c *gin.Context) {
	pullMutex.RLock()
	tag := currentPullTag
	layers := make([]LayerProgress, 0, len(pullLayerProgress))

	totalCurrent := int64(0)
	totalSize := int64(0)
	layerCount := 0
	completedCount := 0

	for id, p := range pullLayerProgress {
		layers = append(layers, LayerProgress{
			ID:       id,
			Status:   p.Status,
			Progress: p.Progress,
			Current:  p.ProgressDetail.Current,
			Total:    p.ProgressDetail.Total,
		})

		if p.ProgressDetail.Total > 0 {
			totalCurrent += p.ProgressDetail.Current
			totalSize += p.ProgressDetail.Total
			layerCount++
			if p.Status == "Pull complete" || p.Status == "Already exists" {
				completedCount++
			}
		}
	}

	percent := 0
	if totalSize > 0 {
		percent = int(float64(totalCurrent) / float64(totalSize) * 100)
	} else if layerCount > 0 {
		percent = int(float64(completedCount) / float64(layerCount) * 100)
	}

	pullMutex.RUnlock()

	response.Success(c, PullProgressResponse{
		Tag:     tag,
		Layers:  layers,
		Percent: percent,
	})
}

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

func (s *ImageService) BuildImage(c *gin.Context) {
	imageName := c.Query("name")
	if imageName == "" {
		response.BadRequest(c, "缺少镜像名称")
		return
	}

	if !validation.ValidateDockerTag(imageName) {
		response.BadRequest(c, "无效的镜像名称")
		return
	}

	dockerfile := c.DefaultQuery("dockerfile", "FROM alpine\nRUN apk add --no-cache curl")
	if dockerfile == "" {
		response.BadRequest(c, "Dockerfile不能为空")
		return
	}

	if len(dockerfile) > 1024*1024 {
		response.BadRequest(c, "Dockerfile内容过长")
		return
	}

	if err := s.client.BuildImage(imageName, dockerfile); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "镜像构建成功", gin.H{"name": imageName})
}

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
