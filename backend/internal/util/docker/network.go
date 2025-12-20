package docker

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// 网络结构体
type Network struct {
	Name       string            `json:"name"`
	ID         string            `json:"id"`
	Created    time.Time         `json:"created"`
	Scope      string            `json:"scope"`
	Driver     string            `json:"driver"`
	EnableIPv6 bool              `json:"enable_ipv6"`
	IPAM       IPAM              `json:"ipam"`
	Internal   bool              `json:"internal"`
	Attachable bool              `json:"attachable"`
	Ingress    bool              `json:"ingress"`
	ConfigFrom string            `json:"config_from"`
	ConfigOnly bool              `json:"config_only"`
	Containers map[string]string `json:"containers"`
	Options    map[string]string `json:"options"`
	Labels     map[string]string `json:"labels"`
}

// IPAM配置
type IPAM struct {
	Driver  string            `json:"driver"`
	Config  []IPAMConfig      `json:"config"`
	Options map[string]string `json:"options"`
}

// IPAM配置项
type IPAMConfig struct {
	Subnet     string            `json:"subnet"`
	IPRange    string            `json:"ip_range"`
	Gateway    string            `json:"gateway"`
	AuxAddress map[string]string `json:"aux_address"`
}

// 获取网络列表
func (c *Client) ListNetworks() ([]*Network, error) {
	ctx := context.Background()
	result, err := c.cli.NetworkList(ctx, client.NetworkListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	networks := make([]*Network, 0, len(result.Items))
	for _, net := range result.Items {
		ipam := IPAM{
			Driver:  net.IPAM.Driver,
			Options: net.IPAM.Options,
		}

		for _, config := range net.IPAM.Config {
			auxAddress := make(map[string]string)
			for k, v := range config.AuxAddress {
				auxAddress[k] = v.String()
			}

			ipam.Config = append(ipam.Config, IPAMConfig{
				Subnet:     config.Subnet.String(),
				IPRange:    config.IPRange.String(),
				Gateway:    config.Gateway.String(),
				AuxAddress: auxAddress,
			})
		}

		networks = append(networks, &Network{
			Name:       net.Name,
			ID:         net.ID,
			Created:    time.Time{},
			Scope:      net.Scope,
			Driver:     net.Driver,
			EnableIPv6: net.EnableIPv6,
			IPAM:       ipam,
			Internal:   net.Internal,
			Attachable: net.Attachable,
			Ingress:    net.Ingress,
			ConfigFrom: string(net.ConfigFrom.Network),
			ConfigOnly: net.ConfigOnly,
			Containers: make(map[string]string),
			Options:    net.Options,
			Labels:     net.Labels,
		})
	}

	return networks, nil
}

// 创建网络
func (c *Client) CreateNetwork(name, driver, subnet, gateway string) (*Network, error) {
	ctx := context.Background()

	var ipamConfig network.IPAMConfig
	if subnet != "" {
		if prefix, err := netip.ParsePrefix(subnet); err == nil {
			ipamConfig.Subnet = prefix
		}
	}
	if gateway != "" {
		if addr, err := netip.ParseAddr(gateway); err == nil {
			ipamConfig.Gateway = addr
		}
	}

	opts := client.NetworkCreateOptions{
		Driver: driver,
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{ipamConfig},
		},
	}

	resp, err := c.cli.NetworkCreate(ctx, name, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	id := resp.ID
	if len(id) > 12 {
		id = id[:12]
	}

	return &Network{
		Name:   name,
		ID:     id,
		Driver: driver,
	}, nil
}

// 删除网络
func (c *Client) RemoveNetwork(id string) error {
	ctx := context.Background()
	_, err := c.cli.NetworkRemove(ctx, id, client.NetworkRemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove network: %w", err)
	}

	return nil
}

// 连接容器到网络
func (c *Client) ConnectContainer(networkID, containerID string, ipv4 string) error {
	ctx := context.Background()

	var endpointConfig *network.EndpointSettings
	if ipv4 != "" {
		addr, err := netip.ParseAddr(ipv4)
		if err != nil {
			return fmt.Errorf("invalid IPv4 address: %w", err)
		}
		endpointConfig = &network.EndpointSettings{
			IPAMConfig: &network.EndpointIPAMConfig{
				IPv4Address: addr,
			},
		}
	}

	_, err := c.cli.NetworkConnect(ctx, networkID, client.NetworkConnectOptions{
		Container:      containerID,
		EndpointConfig: endpointConfig,
	})
	if err != nil {
		return fmt.Errorf("failed to connect container to network: %w", err)
	}

	return nil
}

// 断开容器与网络的连接
func (c *Client) DisconnectContainer(networkID, containerID string) error {
	ctx := context.Background()

	_, err := c.cli.NetworkDisconnect(ctx, networkID, client.NetworkDisconnectOptions{
		Container: containerID,
		Force:     false,
	})
	if err != nil {
		return fmt.Errorf("failed to disconnect container from network: %w", err)
	}

	return nil
}
