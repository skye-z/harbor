package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// 获取网络列表
func (d Docker) GetNetworkList() ([]types.NetworkResource, error) {
	networks, err := d.Session.NetworkList(d.Context, types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	return networks, nil
}

// 获取网络详情
func (d Docker) GetNetworkInfo(id string) (types.NetworkResource, error) {
	info, err := d.Session.NetworkInspect(d.Context, id, types.NetworkInspectOptions{
		Verbose: true,
	})
	if err != nil {
		return types.NetworkResource{}, err
	}

	return info, nil
}

// 创建网络
func (d Docker) CreateNetwork(name string, driver string, ipv6 bool, internal bool, attachable bool, subnet string, gateway string) (string, error) {
	networkResponse, err := d.Session.NetworkCreate(d.Context, name, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         driver,
		EnableIPv6:     ipv6,
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{
				{
					Subnet:  subnet,
					Gateway: gateway,
				},
			},
		},
		Internal:   internal,
		Attachable: attachable,
	})
	if err != nil {
		return "", err
	}

	return networkResponse.ID, nil
}

// 删除网络
func (d Docker) RemoveNetwork(id string) error {
	err := d.Session.NetworkRemove(d.Context, id)
	return err
}

// 容器接入网络
func (d Docker) ConnectNetwork(id string, container string, alias string, ipv4 string, ipv6 string) error {
	err := d.Session.NetworkConnect(d.Context, id, container, &network.EndpointSettings{
		IPAMConfig: &network.EndpointIPAMConfig{
			IPv4Address: ipv4,
			IPv6Address: ipv6,
		},
		Aliases: []string{alias},
	})
	return err
}

// 容器断开网络
func (d Docker) DisconnectNetwork(id string, container string) error {
	err := d.Session.NetworkDisconnect(d.Context, id, container, false)
	return err
}
