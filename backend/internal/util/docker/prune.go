package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

// 清理未使用的容器
func (c *Client) PruneContainers(ctx context.Context) (int64, error) {
	_, err := c.cli.ContainerPrune(ctx, client.ContainerPruneOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to prune containers: %w", err)
	}
	return 0, nil
}

// 清理未使用的镜像
func (c *Client) PruneImages(ctx context.Context) (int64, error) {
	_, err := c.cli.ImagePrune(ctx, client.ImagePruneOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to prune images: %w", err)
	}
	return 0, nil
}

// 清理未使用的卷
func (c *Client) PruneVolumes(ctx context.Context) (int64, error) {
	_, err := c.cli.VolumePrune(ctx, client.VolumePruneOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to prune volumes: %w", err)
	}
	return 0, nil
}

// 清理未使用的网络
func (c *Client) PruneNetworks(ctx context.Context) (int64, error) {
	_, err := c.cli.NetworkPrune(ctx, client.NetworkPruneOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to prune networks: %w", err)
	}
	return 0, nil
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
