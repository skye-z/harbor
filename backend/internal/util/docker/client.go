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

func parseDockerHost(socket string) string {
	if strings.HasPrefix(socket, "unix://") || strings.HasPrefix(socket, "unix:///") {
		return socket
	} else if strings.HasPrefix(socket, "/") {
		return "unix://" + socket
	} else if strings.HasPrefix(socket, "tcp://") || strings.HasPrefix(socket, "http://") {
		return socket
	}
	return socket
}

// 创建新的客户端
func NewClient() (*Client, error) {
	socket := config.GetString("docker.socket")
	host := parseDockerHost(socket)

	cli, err := client.New(client.WithHost(host))
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &Client{
		cli: cli,
	}, nil
}

// 关闭客户端连接
func (c *Client) Close() error {
	if c.cli != nil {
		c.cli.Close()
	}
	return nil
}

// 将容器打包为镜像
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
