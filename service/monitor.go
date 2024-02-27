package service

import (
	"github.com/skye-z/harbor/docker"
	"github.com/skye-z/harbor/monitor"
	"github.com/skye-z/harbor/util"

	"github.com/gin-gonic/gin"
)

type MonintorService struct {
	Client *docker.Docker
}

func NewMonintorService(client *docker.Docker) *MonintorService {
	ms := new(MonintorService)
	return ms
}

// 获取系统信息
func (ms MonintorService) GetDeviceInfo(ctx *gin.Context) {
	info := monitor.GetDeviceInfo()
	util.ReturnData(ctx, true, info)
}

// 获取资源占用情况
func (ms MonintorService) GetUse(ctx *gin.Context) {
	use := monitor.GetUse()
	util.ReturnData(ctx, true, use)
}
