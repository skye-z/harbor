package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Docker struct {
	// 会话上下文
	Context context.Context
	// 守护进程信息
	Config types.Info
	// 会话客户端
	Session *client.Client
}

func NewDocker() (*Docker, error) {
	session, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("[Docker] unable to connect to local docker: %v", err)
	}
	ctx := context.Background()
	info, err := session.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Docker] unable to connect to local docker: %v", err)
	}
	log.Println("[Docker] docker session created")
	return &Docker{
		Context: ctx,
		Config:  info,
		Session: session,
	}, nil
}

func (d Docker) Close() {
	fmt.Println("docker session closed")
	d.Session.Close()
}

// 清理构建缓存
func (d Docker) GetUsage() (types.DiskUsage, error) {
	return d.Session.DiskUsage(d.Context, types.DiskUsageOptions{})
}

// 清理构建缓存
func (d Docker) CleanBuildCache() error {
	_, err := d.Session.BuildCachePrune(d.Context, types.BuildCachePruneOptions{
		All: true,
	})
	return err
}

// 清理未使用的镜像
func (d Docker) CleanImage() error {
	_, err := d.Session.ImagesPrune(d.Context, filters.NewArgs())
	return err
}

// 清理未使用的网络
func (d Docker) CleanNetworks() error {
	_, err := d.Session.NetworksPrune(d.Context, filters.NewArgs())
	return err
}

// 清理未使用的存储卷
func (d Docker) CleanVolumes() error {
	_, err := d.Session.VolumesPrune(d.Context, filters.NewArgs())
	return err
}
