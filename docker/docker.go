package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
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
		return nil, fmt.Errorf("unable to connect to local docker: %v", err)
	}
	ctx := context.Background()
	info, err := session.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to local docker: %v", err)
	}
	fmt.Println("docker session created")
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
