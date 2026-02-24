package service

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/util/docker"
	"github.com/skye-z/harbor/internal/util/response"
)

// 系统管理服务
type DockerService struct {
	client *docker.Client
}

// 创建服务实例
func NewDockerService(client *docker.Client) *DockerService {
	return &DockerService{client: client}
}

// 关闭连接
func (s *DockerService) Close(c *gin.Context) {
	if err := s.client.Close(); err != nil {
		response.Error(c, "关闭连接失败")
		return
	}
	response.SuccessWithMessage(c, "连接已关闭", nil)
}

// 清理未使用的容器
func (s *DockerService) PruneContainers(c *gin.Context) {
	result, err := s.client.PruneContainers(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"reclaimed_space": result})
}

// 清理未使用的镜像
func (s *DockerService) PruneImages(c *gin.Context) {
	result, err := s.client.PruneImages(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"reclaimed_space": result})
}

// 清理未使用的卷
func (s *DockerService) PruneVolumes(c *gin.Context) {
	result, err := s.client.PruneVolumes(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"reclaimed_space": result})
}

// 清理未使用的网络
func (s *DockerService) PruneNetworks(c *gin.Context) {
	result, err := s.client.PruneNetworks(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"reclaimed_space": result})
}

// 清理所有未使用的资源
func (s *DockerService) PruneAll(c *gin.Context) {
	result, err := s.client.PruneAll(c.Request.Context())
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, gin.H{"result": result})
}

// 获取系统信息
func (s *DockerService) GetInfo(c *gin.Context) {
	info, err := s.client.GetInfo()
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	diskInfo, _ := GetDiskUsage()

	data, _ := json.Marshal(info)
	var infoMap map[string]interface{}
	json.Unmarshal(data, &infoMap)

	if diskInfo != nil {
		infoMap["disk_total"] = diskInfo.Total
		infoMap["disk_used"] = diskInfo.Used
		infoMap["disk_available"] = diskInfo.Available
	}

	response.Success(c, infoMap)
}

// 磁盘使用情况
type DiskUsageInfo struct {
	Total     int64 `json:"total"`
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
}

// 获取磁盘使用情况
func GetDiskUsage() (*DiskUsageInfo, error) {
	return &DiskUsageInfo{
		Total:     0,
		Used:      0,
		Available: 0,
	}, nil
}
