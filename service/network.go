package service

import (
	"github.com/skye-z/harbor/docker"
	"github.com/skye-z/harbor/util"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
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

type FormNewtork struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Driver      string `json:"driver"`
	Internal    bool   `json:"internal"`
	Attachable  bool   `json:"attachable"`
	IPv4        bool   `json:"ipv4"`
	IPv4Subnet  string `json:"ipv4Subnet"`
	IPv4Gateway string `json:"ipv4Gateway"`
	IPv4Range   string `json:"ipv4Range"`
	IPv6        bool   `json:"ipv6"`
	IPv6Subnet  string `json:"ipv6Subnet"`
	IPv6Gateway string `json:"ipv6Gateway"`
	IPv6Range   string `json:"ipv6Range"`
}

func (ns NetworkService) Create(ctx *gin.Context) {
	var addObj FormNewtork
	if err := ctx.ShouldBindJSON(&addObj); err != nil {
		util.ReturnMessage(ctx, false, "传入数据无效")
		return
	}

	list := []network.IPAMConfig{}
	if addObj.IPv4 {
		list = append(list, network.IPAMConfig{
			Subnet:  addObj.IPv4Subnet,
			Gateway: addObj.IPv4Gateway,
			IPRange: addObj.IPv4Range,
		})
	}
	if addObj.IPv6 {
		list = append(list, network.IPAMConfig{
			Subnet:  addObj.IPv6Subnet,
			Gateway: addObj.IPv6Gateway,
			IPRange: addObj.IPv6Range,
		})
	}

	id, err := ns.Client.CreateNetwork(addObj.Name, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         addObj.Driver,
		EnableIPv6:     addObj.IPv6,
		IPAM: &network.IPAM{
			Driver: "default",
			Config: list,
		},
		Internal:   addObj.Internal,
		Attachable: addObj.Attachable,
	})
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
