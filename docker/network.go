package docker

import "github.com/docker/docker/api/types"

// 获取网络列表
func (d Docker) GetNetworkList() ([]types.NetworkResource, error) {
	networks, err := d.Session.NetworkList(d.Context, types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	return networks, nil
}
