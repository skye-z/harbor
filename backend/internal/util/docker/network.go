package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/moby/moby/client"
)

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

type IPAM struct {
	Driver  string            `json:"driver"`
	Config  []IPAMConfig      `json:"config"`
	Options map[string]string `json:"options"`
}

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
	_, err := c.cli.NetworkCreate(ctx, name, client.NetworkCreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	return &Network{
		Name:   name,
		ID:     name[:12],
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
