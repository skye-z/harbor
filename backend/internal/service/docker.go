package service

import (
	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/util/docker"
	"github.com/skye-z/harbor/internal/util/response"
)

// DockerService Docker系统管理服务
type DockerService struct {
	client *docker.Client
}

// NewDockerService 创建Docker服务实例
func NewDockerService(client *docker.Client) *DockerService {
	return &DockerService{client: client}
}

// Close 关闭Docker客户端连接
func (s *DockerService) Close(c *gin.Context) {
	if err := s.client.Close(); err != nil {
		response.Error(c, "关闭连接失败")
		return
	}
	response.SuccessWithMessage(c, "连接已关闭", nil)
}

// PruneContainers 清理未使用的容器
func (s *DockerService) PruneContainers(c *gin.Context) {
	result, err := s.client.PruneContainers(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"reclaimed_space": result})
}

// PruneImages 清理未使用的镜像
func (s *DockerService) PruneImages(c *gin.Context) {
	result, err := s.client.PruneImages(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"reclaimed_space": result})
}

// PruneVolumes 清理未使用的卷
func (s *DockerService) PruneVolumes(c *gin.Context) {
	result, err := s.client.PruneVolumes(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"reclaimed_space": result})
}

// PruneNetworks 清理未使用的网络
func (s *DockerService) PruneNetworks(c *gin.Context) {
	result, err := s.client.PruneNetworks(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"reclaimed_space": result})
}

// PruneAll 清理所有未使用的资源
func (s *DockerService) PruneAll(c *gin.Context) {
	result, err := s.client.PruneAll(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"result": result})
}

// GetInfo 获取Docker系统信息
func (s *DockerService) GetInfo(c *gin.Context) {
	info, err := s.client.GetInfo()
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, info)
}
