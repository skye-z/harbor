package docker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/client"
	"github.com/skye-z/harbor/internal/util/config"
)

type Client struct {
	cli     *client.Client
	ctx     context.Context
	cancel  context.CancelFunc
	eventWg sync.WaitGroup
}

type EventCallback func(event events.Message)

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

	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		cli:    cli,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (c *Client) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.eventWg.Wait()
	if c.cli != nil {
		c.cli.Close()
	}
	return nil
}

func (c *Client) StartEventListener(callback EventCallback) error {
	go func() {
		c.eventWg.Add(1)
		defer c.eventWg.Done()

		for {
			select {
			case <-c.ctx.Done():
				log.Println("[Docker] Event listener stopped")
				return
			default:
				err := c.listenEvents(callback)
				if err != nil {
					log.Printf("[Docker] Event listener error: %v, reconnecting in 10s...", err)
					select {
					case <-c.ctx.Done():
						return
					case <-time.After(10 * time.Second):
					}
				}
			}
		}
	}()

	log.Println("[Docker] Event listener started")
	return nil
}

func (c *Client) listenEvents(callback EventCallback) error {
	ctx := context.Background()
	result := c.cli.Events(ctx, client.EventsListOptions{})

	log.Println("[Docker] start listening events for docker")

	reconnect := false

	for !reconnect {
		select {
		case <-c.ctx.Done():
			return nil
		case event := <-result.Messages:
			log.Printf("[Docker] Received event: Type=%s, Action=%s, Actor.ID=%s", event.Type, event.Action, event.Actor.ID)
			if callback != nil {
				callback(event)
			}
		case err := <-result.Err:
			log.Printf("[Docker] Event error: %v", err)
			reconnect = true
		}
	}

	return nil
}

func (c *Client) StopEventListener() {
	if c.cancel != nil {
		c.cancel()
	}
	c.eventWg.Wait()
	log.Println("[Docker] Event listener stopped")
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
