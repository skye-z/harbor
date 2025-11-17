package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/moby/moby/client"
	"github.com/skye-z/harbor/internal/util/config"
)

type Client struct {
	cli *client.Client
}

// 创建新的 Docker 客户端
func NewClient() (*Client, error) {
	socket := config.GetString("docker.socket")

	var host string
	if strings.HasPrefix(socket, "unix://") || strings.HasPrefix(socket, "unix:///") {
		host = socket
	} else if strings.HasPrefix(socket, "/") {
		host = "unix://" + socket
	} else if strings.HasPrefix(socket, "tcp://") || strings.HasPrefix(socket, "http://") {
		host = socket
	} else {
		host = socket
	}

	cli, err := client.New(client.WithHost(host))
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &Client{
		cli: cli,
	}, nil
}

// 关闭 Docker 客户端连接
func (c *Client) Close() error {
	if c.cli != nil {
		c.cli.Close()
	}
	return nil
}

// CommitContainer 将容器打包为镜像
func (c *Client) CommitContainer(containerID, repo, tag string) (string, error) {
	ctx := context.Background()

	resp, err := c.cli.ContainerCommit(ctx, containerID, client.ContainerCommitOptions{
		Reference: fmt.Sprintf("%s:%s", repo, tag),
	})
	if err != nil {
		return "", fmt.Errorf("failed to commit container: %w", err)
	}

	return resp.ID, nil
}
