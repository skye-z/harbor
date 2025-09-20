package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type Container struct {
	ID              string                 `json:"id"`
	Names           []string               `json:"names"`
	Image           string                 `json:"image"`
	ImageID         string                 `json:"image_id"`
	Command         string                 `json:"command"`
	Created         int64                  `json:"created"`
	State           string                 `json:"state"`
	Status          string                 `json:"status"`
	Ports           []Port                 `json:"ports"`
	SizeRw          int64                  `json:"size_rw"`
	SizeRootFs      int64                  `json:"size_root_fs"`
	Labels          map[string]string      `json:"labels"`
	HostConfig      HostConfig             `json:"host_config"`
	NetworkSettings NetworkSettingsSummary `json:"network_settings"`
	Mounts          []Mount                `json:"mounts"`
}

type Port struct {
	IP          string `json:"ip"`
	PrivatePort int    `json:"private_port"`
	PublicPort  int    `json:"public_port"`
	Type        string `json:"type"`
}

type HostConfig struct {
	NetworkMode string            `json:"network_mode"`
	Annotations map[string]string `json:"annotations"`
}

type NetworkSettingsSummary struct {
	Networks map[string]EndpointSettings `json:"networks"`
}

type EndpointSettings struct {
	NetworkID           string   `json:"network_id"`
	EndpointID          string   `json:"endpoint_id"`
	Gateway             string   `json:"gateway"`
	IPAddress           string   `json:"ip_address"`
	IPPrefixLen         int      `json:"ip_prefix_len"`
	IPv6Gateway         string   `json:"ipv6_gateway"`
	GlobalIPv6Address   string   `json:"global_ipv6_address"`
	GlobalIPv6PrefixLen int      `json:"global_ipv6_prefix_len"`
	MacAddress          string   `json:"mac_address"`
	Aliases             []string `json:"aliases"`
}

type Mount struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Driver      string `json:"driver"`
	Mode        string `json:"mode"`
	RW          bool   `json:"rw"`
	Propagation string `json:"propagation"`
}

// 获取容器列表
func (c *Client) ListContainers() ([]*Container, error) {
	ctx := context.Background()
	result, err := c.cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	containers := make([]*Container, 0, len(result.Items))
	for _, cont := range result.Items {
		ports := make([]Port, 0, len(cont.Ports))
		for _, p := range cont.Ports {
			ports = append(ports, Port{
				IP:          p.IP.String(),
				PrivatePort: int(p.PrivatePort),
				PublicPort:  int(p.PublicPort),
				Type:        p.Type,
			})
		}

		container := &Container{
			ID:         cont.ID[:12],
			Names:      cont.Names,
			Image:      cont.Image,
			ImageID:    cont.ImageID,
			Command:    cont.Command,
			Created:    cont.Created,
			State:      string(cont.State),
			Status:     cont.Status,
			Ports:      ports,
			SizeRw:     cont.SizeRw,
			SizeRootFs: cont.SizeRootFs,
			Labels:     cont.Labels,
			HostConfig: HostConfig{
				NetworkMode: string(cont.HostConfig.NetworkMode),
				Annotations: cont.HostConfig.Annotations,
			},
		}

		if cont.NetworkSettings != nil {
			container.NetworkSettings.Networks = make(map[string]EndpointSettings)
			for netName, endpoint := range cont.NetworkSettings.Networks {
				container.NetworkSettings.Networks[netName] = EndpointSettings{
					NetworkID:           endpoint.NetworkID,
					EndpointID:          endpoint.EndpointID,
					Gateway:             endpoint.Gateway.String(),
					IPAddress:           endpoint.IPAddress.String(),
					IPPrefixLen:         endpoint.IPPrefixLen,
					IPv6Gateway:         endpoint.IPv6Gateway.String(),
					GlobalIPv6Address:   endpoint.GlobalIPv6Address.String(),
					GlobalIPv6PrefixLen: endpoint.GlobalIPv6PrefixLen,
					MacAddress:          string(endpoint.MacAddress),
					Aliases:             endpoint.Aliases,
				}
			}
		}

		for _, m := range cont.Mounts {
			container.Mounts = append(container.Mounts, Mount{
				Type:        string(m.Type),
				Name:        m.Name,
				Source:      m.Source,
				Destination: m.Destination,
				Driver:      m.Driver,
				Mode:        m.Mode,
				RW:          m.RW,
				Propagation: string(m.Propagation),
			})
		}

		containers = append(containers, container)
	}

	return containers, nil
}

// 操作容器
func (c *Client) OperationContainer(id string, aciton int) error {
	ctx := context.Background()

	var err error = nil
	switch aciton {
	case 1:
		// 启动容器
		_, err = c.cli.ContainerStart(ctx, id, client.ContainerStartOptions{})
	case 2:
		// 停止容器
		_, err = c.cli.ContainerStop(ctx, id, client.ContainerStopOptions{})
	case 3:
		// 重启容器
		_, err = c.cli.ContainerRestart(ctx, id, client.ContainerRestartOptions{})
	case 4:
		// 删除容器
		_, err = c.cli.ContainerRemove(ctx, id, client.ContainerRemoveOptions{})
	}

	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	return nil
}

// 获取容器日志
func (c *Client) GetContainerLogs(ctx context.Context, id string, tail string) (io.ReadCloser, error) {
	resp, err := c.cli.ContainerLogs(ctx, id, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get container logs: %w", err)
	}

	return resp, nil
}

// 容器统计信息
type ContainerStats struct {
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	NetworkRx float64 `json:"networkRx"`
	NetworkTx float64 `json:"networkTx"`
}

// 获取容器统计信息
func (c *Client) GetContainerStats(id string) (*ContainerStats, error) {
	ctx := context.Background()
	result, err := c.cli.ContainerStats(ctx, id, client.ContainerStatsOptions{Stream: false})
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer result.Body.Close()

	var stats container.StatsResponse
	if err := json.NewDecoder(result.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode stats: %w", err)
	}

	// 计算 CPU 使用率
	var cpuPercent = 0.0
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	onlineCPUs := float64(stats.CPUStats.OnlineCPUs)
	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}

	// 计算内存使用量
	// 兼容不同版本的 Docker API，有些版本可能没有 Stats 字段
	var memoryUsage = float64(stats.MemoryStats.Usage)
	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		memoryUsage = memoryUsage - float64(cache)
	} else if inactiveFile, ok := stats.MemoryStats.Stats["inactive_file"]; ok {
		memoryUsage = memoryUsage - float64(inactiveFile)
	}

	// 计算网络流量
	var networkRx, networkTx float64
	for _, net := range stats.Networks {
		networkRx += float64(net.RxBytes)
		networkTx += float64(net.TxBytes)
	}

	return &ContainerStats{
		CPU:       cpuPercent,
		Memory:    memoryUsage,
		NetworkRx: networkRx,
		NetworkTx: networkTx,
	}, nil
}

// 获取容器进程信息
func (c *Client) GetContainerProcesses(ctx context.Context, id string) (client.ContainerTopResult, error) {
	resp, err := c.cli.ContainerTop(ctx, id, client.ContainerTopOptions{
		Arguments: []string{"pid", "ppid", "user", "cmd"},
	})
	if err != nil {
		return client.ContainerTopResult{}, fmt.Errorf("failed to get container logs: %w", err)
	}

	return resp, nil
}

// 获取容器目录结构
func (c *Client) ListContainerDir(ctx context.Context, id string, path string) (client.ContainerStatPathResult, error) {
	resp, err := c.cli.ContainerStatPath(ctx, id, client.ContainerStatPathOptions{
		Path: path,
	})
	if err != nil {
		return client.ContainerStatPathResult{}, fmt.Errorf("failed to get container logs: %w", err)
	}

	return resp, nil
}
