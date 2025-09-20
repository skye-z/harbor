package docker

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/moby/moby/client"
	"github.com/skye-z/harbor/internal/util/config"
)

type Client struct {
	cli *client.Client
}

// 执行结果
type ExecResult struct {
	Reader io.Reader          // 输出读取器
	Conn   io.ReadWriteCloser // 读写连接
}

// 容器创建配置
type ContainerCreateConfig struct {
	Cmd          []string // 执行命令
	Env          []string // 环境变量
	WorkingDir   string   // 工作目录
	User         string   // 执行用户
	AttachStdin  bool     // 是否附加标准输入
	AttachStdout bool     // 是否附加标准输出
	AttachStderr bool     // 是否附加标准错误输出
	Tty          bool     // 是否使用 TTY（终端）
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

// 在容器内创建执行实例
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

// 附加到执行实例
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

// 调整执行实例的终端大小
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

// 关闭执行实例
func (c *Client) CloseExec(ctx context.Context, execID string) error {
	return nil
}

// 获取容器的详细信息
func (c *Client) GetContainerInfo(ctx context.Context, id string) (*client.ContainerInspectResult, error) {
	inspect, err := c.cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	return &inspect, nil
}

// 获取容器基础信息
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
