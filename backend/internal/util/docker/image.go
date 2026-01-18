package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/moby/moby/client"
)

// 镜像结构体
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
		id := img.ID
		if len(id) > 12 {
			id = id[:12]
		}
		images = append(images, &Image{
			ID:          id,
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

type PullProgress struct {
	Status         string `json:"status"`
	ID             string `json:"id"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
}

// 拉取镜像，支持进度回调
func (c *Client) PullImage(tag string, onProgress func(PullProgress)) error {
	ctx := context.Background()
	resp, err := c.cli.ImagePull(ctx, tag, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()

	decoder := json.NewDecoder(resp)
	for {
		var progress PullProgress
		if err := decoder.Decode(&progress); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode progress: %w", err)
		}
		if onProgress != nil {
			onProgress(progress)
		}
	}

	return nil
}

// 构建镜像
func (c *Client) BuildImage(imageName, dockerfileContent string) error {
	ctx := context.Background()

	tmpDir, err := os.MkdirTemp("", "docker-build")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}

	cmd := exec.CommandContext(ctx, "docker", "build", "-t", imageName, "-f", dockerfilePath, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}

	return nil
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

// 查看镜像详情
func (c *Client) InspectImage(id string) (*client.ImageInspectResult, error) {
	ctx := context.Background()
	result, err := c.cli.ImageInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect image: %w", err)
	}

	return &result, nil
}

// 打标签
func (c *Client) TagImage(imageID, tag string) error {
	ctx := context.Background()
	_, err := c.cli.ImageTag(ctx, client.ImageTagOptions{
		Source: imageID,
		Target: tag,
	})
	if err != nil {
		return fmt.Errorf("failed to tag image: %w", err)
	}
	return nil
}

// 推送镜像
func (c *Client) PushImage(tag string) error {
	ctx := context.Background()
	resp, err := c.cli.ImagePush(ctx, tag, client.ImagePushOptions{})
	if err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}
	defer resp.Close()

	_, err = io.Copy(io.Discard, resp)
	if err != nil {
		return fmt.Errorf("failed to read push output: %w", err)
	}

	return nil
}

// 搜索镜像
type SearchResult struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StarCount   int    `json:"star_count"`
	IsOfficial  bool   `json:"is_official"`
	IsAutomated bool   `json:"is_automated"`
}

func (c *Client) SearchImages(query string, limit int) ([]SearchResult, error) {
	ctx := context.Background()
	result, err := c.cli.ImageSearch(ctx, query, client.ImageSearchOptions{
		Limit: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search images: %w", err)
	}

	results := make([]SearchResult, 0, len(result.Items))
	for _, r := range result.Items {
		results = append(results, SearchResult{
			Name:        r.Name,
			Description: r.Description,
			StarCount:   r.StarCount,
			IsOfficial:  r.IsOfficial,
			IsAutomated: r.IsAutomated,
		})
	}

	return results, nil
}
