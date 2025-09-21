package docker

import (
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
