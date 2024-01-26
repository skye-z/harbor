package service

import (
	"harbor/docker"
	"harbor/util"

	"github.com/gin-gonic/gin"
)

type NetworkService struct {
	Client *docker.Docker
}

func NewNetworkService(client *docker.Docker) *NetworkService {
	ns := new(NetworkService)
	ns.Client = client
	return ns
}

func (ns NetworkService) GetList(ctx *gin.Context) {
	list, err := ns.Client.GetNetworkList()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

func (ns NetworkService) GetInfo(ctx *gin.Context) {
	id := ctx.Query("id")
	info, err := ns.Client.GetNetworkInfo(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取网络详情失败")
	} else {
		util.ReturnData(ctx, true, info)
	}
}

func (ns NetworkService) Create(ctx *gin.Context) {
	name := ctx.Query("name")
	driver := ctx.Query("driver")
	ipv6 := ctx.Query("ipv6")
	internal := ctx.Query("internal")
	attachable := ctx.Query("attachable")
	subnet := ctx.Query("subnet")
	gateway := ctx.Query("gateway")
	id, err := ns.Client.CreateNetwork(name, driver, ipv6 == "1", internal == "1", attachable == "1", subnet, gateway)
	if err != nil {
		util.ReturnMessage(ctx, false, "创建网络失败")
	} else {
		util.ReturnData(ctx, true, id)
	}
}

func (ns NetworkService) Remove(ctx *gin.Context) {
	id := ctx.Query("id")
	err := ns.Client.RemoveNetwork(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "删除网络失败")
	} else {
		util.ReturnMessage(ctx, true, "删除网络成功")
	}
}

func (ns NetworkService) Connect(ctx *gin.Context) {
	id := ctx.Query("id")
	container := ctx.Query("container")
	alias := ctx.Query("alias")
	ipv4 := ctx.Query("ipv4")
	ipv6 := ctx.Query("ipv6")
	err := ns.Client.ConnectNetwork(id, container, alias, ipv4, ipv6)
	if err != nil {
		util.ReturnMessageData(ctx, false, "接入网络失败", err.Error())
	} else {
		util.ReturnMessage(ctx, true, "接入网络成功")
	}
}

func (ns NetworkService) Disconnect(ctx *gin.Context) {
	id := ctx.Query("id")
	container := ctx.Query("container")
	err := ns.Client.DisconnectNetwork(id, container)
	if err != nil {
		util.ReturnMessage(ctx, false, "断开网络失败")
	} else {
		util.ReturnMessage(ctx, true, "断开网络成功")
	}
}
