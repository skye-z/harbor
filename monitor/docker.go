package monitor

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

var (
	errorLogged   = false
	errorLoggedMu sync.Mutex
)

func ListenDockerEvents() {
	for {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			log.Println("[Monitor] listening service startup failed")
		}

		ctx, cancel := context.WithCancel(context.Background())
		eventChan, errChan := cli.Events(ctx, types.EventsOptions{})

		log.Println("[Monitor] start listening events for docker")

		reconnect := false

		for !reconnect {
			select {
			case event := <-eventChan:
				handleEvent(event)
			case err := <-errChan:
				handleError(err, cli, cancel)
				reconnect = true
			}
		}
		cli.Close()
		cancel()
		time.Sleep(10 * time.Second)
	}
}

func handleEvent(event events.Message) {
	if event.Type == "container" {
		if event.Action == "start" {
			// 容器启动
			log.Println("[Monitor] docker container", event.ID, "start")
		} else if event.Action == "die" {
			// 容器停止
			log.Println("[Monitor] docker container", event.ID, "die")
		}
	} else if event.Type == "daemon" {
		if event.Action == "shutdown" {
			// Docker 服务停止
			log.Println("[Monitor] docker daemon shutdown")
		}
	}
}

func handleError(err error, cli *client.Client, cancel context.CancelFunc) {
	errorLoggedMu.Lock()
	defer errorLoggedMu.Unlock()

	if !errorLogged {
		log.Println("[Monitor] docker connection disconnected")
		errorLogged = true
		cli.Close()
		cancel()
		log.Println("[Monitor] stopped docker event listener.")
	}
}
