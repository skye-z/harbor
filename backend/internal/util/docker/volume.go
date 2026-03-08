package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/moby/moby/client"
)

type Volume struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	CreatedAt  time.Time         `json:"created_at"`
	Status     map[string]string `json:"status"`
	Labels     map[string]string `json:"labels"`
	Scope      string            `json:"scope"`
	Options    map[string]string `json:"options"`
	UsageData  *UsageData        `json:"usage_data"`
}

type UsageData struct {
	RefCount int64 `json:"ref_count"`
	Size     int64 `json:"size"`
}

// 获取卷列表
func (c *Client) ListVolumes() ([]*Volume, error) {
	ctx := context.Background()
	result, err := c.cli.VolumeList(ctx, client.VolumeListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	volumes := make([]*Volume, 0, len(result.Items))
	for _, vol := range result.Items {
		var usageData *UsageData
		if vol.UsageData != nil {
			usageData = &UsageData{
				RefCount: vol.UsageData.RefCount,
				Size:     vol.UsageData.Size,
			}
		}

		status := make(map[string]string)
		for k, v := range vol.Status {
			status[k] = fmt.Sprintf("%v", v)
		}

		createdAt, _ := time.Parse(time.RFC3339, vol.CreatedAt)

		volumes = append(volumes, &Volume{
			ID:         vol.Name, // Docker卷使用名称作为ID
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
			CreatedAt:  createdAt,
			Status:     status,
			Labels:     vol.Labels,
			Scope:      vol.Scope,
			Options:    vol.Options,
			UsageData:  usageData,
		})
	}

	return volumes, nil
}

// 创建卷
func (c *Client) CreateVolume(name, driver string) (*Volume, error) {
	ctx := context.Background()
	options := client.VolumeCreateOptions{
		Name:   name,
		Driver: driver,
	}

	_, err := c.cli.VolumeCreate(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume: %w", err)
	}

	return &Volume{
		Name:   name,
		Driver: driver,
	}, nil
}

// 删除卷
func (c *Client) RemoveVolume(id string) error {
	ctx := context.Background()
	_, err := c.cli.VolumeRemove(ctx, id, client.VolumeRemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove volume: %w", err)
	}

	return nil
}
