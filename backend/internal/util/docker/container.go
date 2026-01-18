package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

type ContainerStats struct {
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	NetworkRx float64 `json:"networkRx"`
	NetworkTx float64 `json:"networkTx"`
}

type ExecResult struct {
	Reader io.Reader
	Conn   io.ReadWriteCloser
}

type ContainerCreateConfig struct {
	Cmd          []string
	Env          []string
	WorkingDir   string
	User         string
	AttachStdin  bool
	AttachStdout bool
	AttachStderr bool
	Tty          bool
}

func (c *Client) ListContainers() ([]*Container, error) {
	ctx := context.Background()
	result, err := c.cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	containers := make([]*Container, 0, len(result.Items))
	for _, cont := range result.Items {
		id := cont.ID
		if len(id) > 12 {
			id = id[:12]
		}

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
			ID:         id,
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

func (c *Client) GetContainerInfo(ctx context.Context, id string) (*client.ContainerInspectResult, error) {
	inspect, err := c.cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	return &inspect, nil
}

func (c *Client) GetContainerDetails(ctx context.Context, id string) (map[string]interface{}, error) {
	inspect, err := c.cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	return map[string]any{
		"id":    inspect.Container.ID,
		"name":  strings.TrimPrefix(inspect.Container.Name, "/"),
		"image": inspect.Container.Config.Image,
	}, nil
}

func (c *Client) OperationContainer(id string, action int) error {
	ctx := context.Background()

	var err error = nil
	switch action {
	case 1:
		_, err = c.cli.ContainerStart(ctx, id, client.ContainerStartOptions{})
	case 2:
		_, err = c.cli.ContainerStop(ctx, id, client.ContainerStopOptions{})
	case 3:
		_, err = c.cli.ContainerRestart(ctx, id, client.ContainerRestartOptions{})
	case 4:
		_, err = c.cli.ContainerRemove(ctx, id, client.ContainerRemoveOptions{})
	case 5:
		_, err = c.cli.ContainerPause(ctx, id, client.ContainerPauseOptions{})
	case 6:
		_, err = c.cli.ContainerUnpause(ctx, id, client.ContainerUnpauseOptions{})
	}

	if err != nil {
		return fmt.Errorf("failed to operate container: %w", err)
	}
	return nil
}

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

	var memoryUsage = float64(stats.MemoryStats.Usage)
	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		memoryUsage = memoryUsage - float64(cache)
	} else if inactiveFile, ok := stats.MemoryStats.Stats["inactive_file"]; ok {
		memoryUsage = memoryUsage - float64(inactiveFile)
	}

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

func (c *Client) GetContainerProcesses(ctx context.Context, id string) (client.ContainerTopResult, error) {
	resp, err := c.cli.ContainerTop(ctx, id, client.ContainerTopOptions{
		Arguments: []string{},
	})
	if err != nil {
		return client.ContainerTopResult{}, fmt.Errorf("failed to get container processes: %w", err)
	}

	return resp, nil
}

func (c *Client) ListContainerDir(ctx context.Context, id string, path string) (client.ContainerStatPathResult, error) {
	resp, err := c.cli.ContainerStatPath(ctx, id, client.ContainerStatPathOptions{
		Path: path,
	})
	if err != nil {
		return client.ContainerStatPathResult{}, fmt.Errorf("failed to get container logs: %w", err)
	}

	return resp, nil
}

func (c *Client) CopyFromContainer(ctx context.Context, containerID, srcPath, dstPath string) error {
	if containsPathTraversal(srcPath) {
		return fmt.Errorf("invalid source path: potential path traversal")
	}

	result, err := c.cli.CopyFromContainer(ctx, containerID, client.CopyFromContainerOptions{
		SourcePath: srcPath,
	})
	if err != nil {
		return fmt.Errorf("failed to copy from container: %w", err)
	}
	defer result.Content.Close()

	dstDir := filepath.Join(dstPath, containerID)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	baseName := filepath.Base(srcPath)
	if containsPathTraversal(baseName) {
		return fmt.Errorf("invalid filename: potential path traversal")
	}

	dstFile := filepath.Join(dstDir, baseName)
	file, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, result.Content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func containsPathTraversal(path string) bool {
	cleaned := filepath.Clean(path)
	return cleaned != path ||
		path == ".." ||
		len(path) >= 3 && (path[:3] == "../" || path[len(path)-3:] == "/..") ||
		len(path) >= 4 && (path[:4] == "/../" || path[len(path)-4:] == "/..")
}

func (c *Client) CopyToContainer(ctx context.Context, containerID, srcPath, dstPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	_, err = c.cli.CopyToContainer(ctx, containerID, client.CopyToContainerOptions{
		DestinationPath:           dstPath,
		Content:                   strings.NewReader(string(content)),
		AllowOverwriteDirWithFile: true,
	})
	if err != nil {
		return fmt.Errorf("failed to copy to container: %w", err)
	}

	return nil
}

func (c *Client) CreateExec(ctx context.Context, containerID string, config *ContainerCreateConfig) (string, error) {
	execConfig := client.ExecCreateOptions{
		Cmd:          config.Cmd,
		Env:          config.Env,
		WorkingDir:   config.WorkingDir,
		User:         config.User,
		AttachStdin:  config.AttachStdin,
		AttachStdout: config.AttachStdout,
		AttachStderr: config.AttachStderr,
		TTY:          config.Tty,
	}

	resp, err := c.cli.ExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec: %w", err)
	}

	return resp.ID, nil
}

func (c *Client) AttachToExec(ctx context.Context, execID string, detach bool) (*ExecResult, error) {
	hijackedResp, err := c.cli.ExecAttach(ctx, execID, client.ExecAttachOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to attach to exec: %w", err)
	}

	return &ExecResult{
		Reader: hijackedResp.Reader,
		Conn:   hijackedResp.Conn,
	}, nil
}

func (c *Client) ExecResize(ctx context.Context, execID string, rows, cols int) error {
	_, err := c.cli.ExecResize(ctx, execID, client.ExecResizeOptions{
		Height: uint(rows),
		Width:  uint(cols),
	})
	if err != nil {
		return fmt.Errorf("failed to resize exec: %w", err)
	}
	return nil
}

func (c *Client) CloseExec(ctx context.Context, execID string) error {
	_, err := c.cli.ExecInspect(ctx, execID, client.ExecInspectOptions{})
	if err != nil {
		return fmt.Errorf("failed to inspect exec: %w", err)
	}
	return nil
}

func (c *Client) CreateContainer(imageName, containerName string, cmd []string, env []string, workingDir string, hostConfigData interface{}, exposedPortsData map[string]interface{}, tty, openStdin, autoRemove bool) (*Container, error) {
	ctx := context.Background()

	config := &container.Config{
		Image:      imageName,
		Cmd:        cmd,
		Env:        env,
		WorkingDir: workingDir,
		Tty:        tty,
		OpenStdin:  openStdin,
	}

	hostConfig := &container.HostConfig{
		AutoRemove: autoRemove,
	}

	// 处理ExposedPorts
	if exposedPortsData != nil && len(exposedPortsData) > 0 {
		config.ExposedPorts = make(map[string]struct{})
		for port := range exposedPortsData {
			config.ExposedPorts[port] = struct{}{}
		}
	}

	// 手动处理hostConfig数据
	if hostConfigData != nil {
		if hc, ok := hostConfigData.(map[string]interface{}); ok {
			// NetworkMode
			if nm, ok := hc["NetworkMode"].(string); ok && nm != "" {
				hostConfig.NetworkMode = container.NetworkMode(nm)
			}

			// RestartPolicy
			if rp, ok := hc["RestartPolicy"].(map[string]interface{}); ok {
				if name, ok := rp["Name"].(string); ok {
					hostConfig.RestartPolicy.Name = name
				}
			}

			// Memory
			if mem, ok := hc["Memory"].(float64); ok && mem > 0 {
				hostConfig.Memory = int64(mem)
			}

			// CPUShares
			if cpu, ok := hc["CPUShares"].(float64); ok && cpu > 0 {
				hostConfig.CPUShares = int64(cpu)
			}

			// Binds
			if binds, ok := hc["Binds"].([]interface{}); ok {
				for _, b := range binds {
					if bind, ok := b.(string); ok {
						hostConfig.Binds = append(hostConfig.Binds, bind)
					}
				}
			}

			// PortBindings
			if pb, ok := hc["PortBindings"].(map[string]interface{}); ok {
				hostConfig.PortBindings = make(map[string][]struct{ HostPort string })
				for port, bindings := range pb {
					if bindingList, ok := bindings.([]interface{}); ok {
						var ports []struct{ HostPort string }
						for _, b := range bindingList {
							if binding, ok := b.(map[string]interface{}); ok {
								if hp, ok := binding["HostPort"].(string); ok {
									ports = append(ports, struct{ HostPort string }{HostPort: hp})
								}
							}
						}
						hostConfig.PortBindings[port] = ports
					}
				}
			}
		}
	}

	resp, err := c.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:     config,
		HostConfig: hostConfig,
		Name:       containerName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	id := resp.ID
	if len(id) > 12 {
		id = id[:12]
	}

	return &Container{
		ID:    id,
		Image: imageName,
	}, nil
}

func (c *Client) RenameContainer(containerID, newName string) error {
	ctx := context.Background()
	_, err := c.cli.ContainerRename(ctx, containerID, client.ContainerRenameOptions{
		NewName: newName,
	})
	if err != nil {
		return fmt.Errorf("failed to rename container: %w", err)
	}
	return nil
}
