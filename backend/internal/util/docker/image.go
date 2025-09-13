package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/moby/moby/client"
)

type Image struct {
	ID          string            `json:"id"`
	ParentID    string            `json:"parent_id"`
	RepoTags    []string          `json:"repo_tags"`
	RepoDigests []string          `json:"repo_digests"`
	Created     int64             `json:"created"`
	Size        int64             `json:"size"`
	VirtualSize int64             `json:"virtual_size"`
	SharedSize  int64             `json:"shared_size"`
	Labels      map[string]string `json:"labels"`
	Containers  int64             `json:"containers"`
}

// 获取镜像列表
func (c *Client) ListImages() ([]*Image, error) {
	ctx := context.Background()
	result, err := c.cli.ImageList(ctx, client.ImageListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	images := make([]*Image, 0, len(result.Items))
	for _, img := range result.Items {
		images = append(images, &Image{
			ID:          img.ID[:12],
			ParentID:    img.ParentID,
			RepoTags:    img.RepoTags,
			RepoDigests: img.RepoDigests,
			Created:     img.Created,
			Size:        img.Size,
			SharedSize:  img.SharedSize,
			Labels:      img.Labels,
			Containers:  img.Containers,
		})
	}

	return images, nil
}

// 拉取镜像
func (c *Client) PullImage(tag string) error {
	ctx := context.Background()
	resp, err := c.cli.ImagePull(ctx, tag, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()

	_, err = io.Copy(io.Discard, resp)
	if err != nil {
		return fmt.Errorf("failed to read pull output: %w", err)
	}

	return nil
}

// 构建镜像
func (c *Client) BuildImage(imageName, dockerfileContent string) error {
	return fmt.Errorf("none")
}

// 删除镜像
func (c *Client) RemoveImage(id string) error {
	ctx := context.Background()
	_, err := c.cli.ImageRemove(ctx, id, client.ImageRemoveOptions{Force: true})
	if err != nil {
		return fmt.Errorf("failed to remove image: %w", err)
	}
	return nil
}

// 获取镜像详情
func (c *Client) InspectImage(id string) (*client.ImageInspectResult, error) {
	ctx := context.Background()
	result, err := c.cli.ImageInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect image: %w", err)
	}

	return &result, nil
}
