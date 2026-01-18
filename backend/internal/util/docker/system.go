package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"syscall"

	"github.com/moby/moby/client"
)

type SystemInfo struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Containers        int           `json:"containers"`
	ContainersRunning int           `json:"containers_running"`
	ContainersPaused  int           `json:"containers_paused"`
	ContainersStopped int           `json:"containers_stopped"`
	Images            int           `json:"images"`
	Driver            string        `json:"driver"`
	DriverStatus      []interface{} `json:"driver_status"`
	DockerRootDir     string        `json:"docker_root_dir"`
	LoggingDriver     string        `json:"logging_driver"`
	CgroupDriver      string        `json:"cgroup_driver"`
	CgroupVersion     string        `json:"cgroup_version"`
	NEventsListener   int           `json:"n_events_listener"`
	KernelVersion     string        `json:"kernel_version"`
	OperatingSystem   string        `json:"operating_system"`
	OSType            string        `json:"os_type"`
	Architecture      string        `json:"architecture"`
	NCPU              int           `json:"n_cpu"`
	MemTotal          int64         `json:"mem_total"`
	IPv4Forwarding    bool          `json:"ipv4_forwarding"`
	Debug             bool          `json:"debug"`
	NFd               int           `json:"n_fd"`
	OomKillDisable    bool          `json:"oom_kill_disable"`
	ServerVersion     string        `json:"server_version"`
	SecurityOptions   []string      `json:"security_options"`
	Labels            []string      `json:"labels"`
	Version           VersionInfo   `json:"version"`
	DiskTotal         int64         `json:"disk_total"`
	DiskUsed          int64         `json:"disk_used"`
	DiskAvailable     int64         `json:"disk_available"`
}

type VersionInfo struct {
	Platform struct {
		Name string `json:"name"`
	} `json:"Platform"`
	Components []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Details any    `json:"details"`
	} `json:"Components"`
	Version       string `json:"Version"`
	GitCommit     string `json:"GitCommit"`
	GoVersion     string `json:"GoVersion"`
	Os            string `json:"Os"`
	Arch          string `json:"Arch"`
	KernelVersion string `json:"KernelVersion"`
	Experimental  bool   `json:"Experimental"`
	BuildTime     string `json:"BuildTime"`
}

// 检测 Docker 状态
func (c *Client) Ping() error {
	ctx := context.Background()
	_, err := c.cli.Ping(ctx, client.PingOptions{})
	if err != nil {
		return fmt.Errorf("docker not available: %w", err)
	}
	return nil
}

// 获取 Docker 信息
func (c *Client) GetInfo() (*SystemInfo, error) {
	ctx := context.Background()
	result, err := c.cli.Info(ctx, client.InfoOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get info: %w", err)
	}

	driverStatus := make([]interface{}, len(result.Info.DriverStatus))
	for i, status := range result.Info.DriverStatus {
		driverStatus[i] = status
	}

	versionResult, err := c.cli.ServerVersion(ctx, client.ServerVersionOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get docker version: %w", err)
	}

	data, _ := json.Marshal(versionResult)
	var version VersionInfo
	json.Unmarshal(data, &version)

	diskInfo, _ := GetDiskUsage()

	return &SystemInfo{
		ID:                result.Info.ID,
		Name:              result.Info.Name,
		Containers:        result.Info.Containers,
		ContainersRunning: result.Info.ContainersRunning,
		ContainersPaused:  result.Info.ContainersPaused,
		ContainersStopped: result.Info.ContainersStopped,
		Images:            result.Info.Images,
		Driver:            result.Info.Driver,
		DriverStatus:      driverStatus,
		DockerRootDir:     result.Info.DockerRootDir,
		LoggingDriver:     result.Info.LoggingDriver,
		CgroupDriver:      result.Info.CgroupDriver,
		CgroupVersion:     result.Info.CgroupVersion,
		NEventsListener:   result.Info.NEventsListener,
		KernelVersion:     result.Info.KernelVersion,
		OperatingSystem:   result.Info.OperatingSystem,
		OSType:            result.Info.OSType,
		Architecture:      result.Info.Architecture,
		NCPU:              result.Info.NCPU,
		MemTotal:          result.Info.MemTotal,
		IPv4Forwarding:    result.Info.IPv4Forwarding,
		Debug:             result.Info.Debug,
		NFd:               result.Info.NFd,
		OomKillDisable:    result.Info.OomKillDisable,
		ServerVersion:     result.Info.ServerVersion,
		SecurityOptions:   result.Info.SecurityOptions,
		Labels:            result.Info.Labels,
		Version:           version,
		DiskTotal:         diskInfo.Total,
		DiskUsed:          diskInfo.Used,
		DiskAvailable:     diskInfo.Available,
	}, nil
}

type DiskUsageInfo struct {
	Total     int64 `json:"total"`
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
}

func GetDiskUsage() (*DiskUsageInfo, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs("/", &fs)
	if err != nil {
		return nil, err
	}

	total := fs.Blocks * uint64(fs.Bsize)
	available := fs.Bavail * uint64(fs.Bsize)
	used := (fs.Blocks - fs.Bfree) * uint64(fs.Bsize)

	return &DiskUsageInfo{
		Total:     int64(total),
		Used:      int64(used),
		Available: int64(available),
	}, nil
}

// 清理未使用的容器
func (c *Client) PruneContainers(ctx context.Context) (int64, error) {
	result, err := c.cli.ContainerPrune(ctx, client.ContainerPruneOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to prune containers: %w", err)
	}
	return int64(result.Report.SpaceReclaimed), nil
}

// 清理未使用的镜像
func (c *Client) PruneImages(ctx context.Context) (int64, error) {
	result, err := c.cli.ImagePrune(ctx, client.ImagePruneOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to prune images: %w", err)
	}
	return int64(result.Report.SpaceReclaimed), nil
}

// 清理未使用的卷
func (c *Client) PruneVolumes(ctx context.Context) (int64, error) {
	result, err := c.cli.VolumePrune(ctx, client.VolumePruneOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to prune volumes: %w", err)
	}
	return int64(result.Report.SpaceReclaimed), nil
}

// 清理未使用的网络
func (c *Client) PruneNetworks(ctx context.Context) (int64, error) {
	_, err := c.cli.NetworkPrune(ctx, client.NetworkPruneOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to prune networks: %w", err)
	}
	return 1, nil
}

// 清理所有未使用的资源
func (c *Client) PruneAll(ctx context.Context) (map[string]int64, error) {
	result := make(map[string]int64)

	containerSpace, err := c.PruneContainers(ctx)
	if err != nil {
		return nil, err
	}
	result["containers"] = containerSpace

	imageSpace, err := c.PruneImages(ctx)
	if err != nil {
		return nil, err
	}
	result["images"] = imageSpace

	volumeSpace, err := c.PruneVolumes(ctx)
	if err != nil {
		return nil, err
	}
	result["volumes"] = volumeSpace

	networkSpace, err := c.PruneNetworks(ctx)
	if err != nil {
		return nil, err
	}
	result["networks"] = networkSpace

	return result, nil
}
