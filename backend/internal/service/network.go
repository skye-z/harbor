package service

import (
	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/util/docker"
	"github.com/skye-z/harbor/internal/util/response"
)

// 网络管理服务
type NetworkService struct {
	client *docker.Client
}

// 创建网络服务实例
func NewNetworkService(client *docker.Client) *NetworkService {
	return &NetworkService{client: client}
}

// 获取网络列表
func (s *NetworkService) GetList(c *gin.Context) {
	networks, err := s.client.ListNetworks()
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, networks)
}

// 创建网络
func (s *NetworkService) Create(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		response.BadRequest(c, "缺少网络名称")
		return
	}

	driver := c.DefaultQuery("driver", "bridge")

	network, err := s.client.CreateNetwork(name, driver, "", "")
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "网络创建成功", network)
}

// 删除网络
func (s *NetworkService) Remove(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少网络ID")
		return
	}

	if err := s.client.RemoveNetwork(id); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "网络删除成功", gin.H{"id": id})
}

// 连接容器到网络
func (s *NetworkService) ConnectContainer(c *gin.Context) {
	networkID := c.Query("network_id")
	if networkID == "" {
		response.BadRequest(c, "缺少网络ID")
		return
	}

	containerID := c.Query("container_id")
	if containerID == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	ipv4 := c.DefaultQuery("ipv4", "")

	if err := s.client.ConnectContainer(networkID, containerID, ipv4); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "容器已连接到网络", gin.H{
		"network_id":   networkID,
		"container_id": containerID,
	})
}

// 断开容器与网络的连接
func (s *NetworkService) DisconnectContainer(c *gin.Context) {
	networkID := c.Query("network_id")
	if networkID == "" {
		response.BadRequest(c, "缺少网络ID")
		return
	}

	containerID := c.Query("container_id")
	if containerID == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	if err := s.client.DisconnectContainer(networkID, containerID); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "容器已断开与网络的连接", gin.H{
		"network_id":   networkID,
		"container_id": containerID,
	})
}
