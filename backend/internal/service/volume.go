package service

import (
	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/util/docker"
	"github.com/skye-z/harbor/internal/util/response"
)

// 卷管理服务
type VolumeService struct {
	client *docker.Client
}

// 创建卷服务实例
func NewVolumeService(client *docker.Client) *VolumeService {
	return &VolumeService{client: client}
}

// 获取卷列表
func (s *VolumeService) GetList(c *gin.Context) {
	volumes, err := s.client.ListVolumes()
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, volumes)
}

// 创建卷
func (s *VolumeService) Create(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		response.BadRequest(c, "缺少卷名称")
		return
	}

	driver := c.DefaultQuery("driver", "local")

	volume, err := s.client.CreateVolume(name, driver)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "卷创建成功", volume)
}

// 删除卷
func (s *VolumeService) Remove(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少卷ID")
		return
	}

	if err := s.client.RemoveVolume(id); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "卷删除成功", gin.H{"id": id})
}
