package docker

import (
	"encoding/json"
	"harbor/util"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/gorilla/websocket"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// 获取容器列表
func (d Docker) GetContainerList() ([]types.Container, error) {
	containers, err := d.Session.ContainerList(d.Context, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// 获取容器详情
func (d Docker) GetContainerInfo(id string) (types.ContainerJSON, error) {
	containerInfo, err := d.Session.ContainerInspect(d.Context, id)
	if err != nil {
		return types.ContainerJSON{}, err
	}

	return containerInfo, nil
}

// 获取容器原始详情
func (d Docker) GetContainerInfoWithRaw(id string) ([]byte, error) {
	_, json, err := d.Session.ContainerInspectWithRaw(d.Context, id, true)
	if err != nil {
		return nil, err
	}

	return json, nil
}

// 获取容器文件变动
func (d Docker) GetContainerDiff(id string) ([]container.FilesystemChange, error) {
	containerDiff, err := d.Session.ContainerDiff(d.Context, id)
	if err != nil {
		return nil, err
	}

	return containerDiff, nil
}

// 获取容器统计信息
func (d Docker) GetContainerStat(id string) (*types.Stats, error) {
	stats, err := d.Session.ContainerStatsOneShot(d.Context, id)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	var stat types.Stats
	if err := json.NewDecoder(stats.Body).Decode(&stat); err != nil {
		return nil, err
	}

	return &stat, nil
}

// 获取容器统计信息
func (d Docker) GetContainerLogs(id string, tail string) ([]string, error) {
	logs, err := d.Session.ContainerLogs(d.Context, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
		Tail:       tail,
	})
	if err != nil {
		return nil, err
	}
	defer logs.Close()
	cleanedLogs := util.ProcessLogs(logs)

	return cleanedLogs, nil
}

// 创建容器
func (d Docker) CreateContainer(imageName string, containerName string, command []string) (string, error) {
	// 创建容器配置
	containerConfig := &container.Config{
		Image:        imageName,
		Cmd:          command,
		Tty:          false,
		AttachStdout: true,
		AttachStderr: true,
		Env:          []string{"KEY=VALUE"}, // 可以设置环境变量
	}

	hostConfig := &container.HostConfig{}
	networkingConfig := &network.NetworkingConfig{}
	platform := &v1.Platform{}
	// 创建容器
	resp, err := d.Session.ContainerCreate(d.Context, containerConfig, hostConfig, networkingConfig, platform, containerName)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// 启动容器
func (d Docker) StartContainer(id string) error {
	err := d.Session.ContainerStart(d.Context, id, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

// 重启容器(超时强制停止)
func (d Docker) RestartContainer(id string, timeout *int) error {
	err := d.Session.ContainerRestart(d.Context, id, container.StopOptions{Timeout: timeout})
	if err != nil {
		return err
	}

	return nil
}

// 停止容器(超时强制停止)
func (d Docker) StopContainer(id string, timeout *int) error {
	err := d.Session.ContainerStop(d.Context, id, container.StopOptions{Timeout: timeout})
	if err != nil {
		return err
	}

	return nil
}

// 关闭容器(发送信号通知容器关闭)
func (d Docker) KillContainer(id string) error {
	err := d.Session.ContainerKill(d.Context, id, "SIGTERM")
	if err != nil {
		return err
	}

	return nil
}

// 挂起容器
func (d Docker) PauseContainer(id string) error {
	err := d.Session.ContainerPause(d.Context, id)
	if err != nil {
		return err
	}

	return nil
}

// 从挂起中恢复容器
func (d Docker) UnpauseContainer(id string) error {
	err := d.Session.ContainerUnpause(d.Context, id)
	if err != nil {
		return err
	}

	return nil
}

// 删除容器
func (d Docker) RemoveContainer(id string, removeVolumes bool, removeLinks bool, force bool) error {
	err := d.Session.ContainerRemove(d.Context, id, types.ContainerRemoveOptions{
		RemoveVolumes: removeVolumes,
		RemoveLinks:   removeLinks,
		Force:         force,
	})
	if err != nil {
		return err
	}

	return nil
}

// 重命名容器
func (d Docker) RenameContainer(id string, name string) error {
	err := d.Session.ContainerRename(d.Context, id, name)
	if err != nil {
		return err
	}

	return nil
}

// 获取容器内部进程信息
func (d Docker) GetContainerProcesses(id string) ([][]string, error) {
	processes, err := d.Session.ContainerTop(d.Context, id, nil)
	if err != nil {
		return nil, err
	}

	return processes.Processes, nil
}

// 创建终端
func (d Docker) CreateTerminal(conn *websocket.Conn, containerID string, cmd string, cols uint, rows uint) error {
	createResp, err := d.Session.ContainerExecCreate(d.Context, containerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{cmd},
	})
	if err != nil {
		return err
	}

	log.Println("[Docker] create terminal")
	attachResp, err := d.Session.ContainerExecAttach(d.Context, createResp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		return err
	}

	log.Println("[Docker] connect terminal")
	err = d.Session.ContainerExecResize(d.Context, createResp.ID, types.ResizeOptions{
		Height: rows,
		Width:  cols,
	})
	if err != nil {
		return err
	}

	// 创建一个退出通道，用于协程之间的通知
	exitChan := make(chan struct{})

	go func() {
		defer close(exitChan)
		for {
			buffer := make([]byte, 4096)
			_, err := attachResp.Reader.Read(buffer)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					defer attachResp.Close()
				}
				break
			}
			conn.WriteMessage(websocket.TextMessage, buffer)
		}
	}()

	go func() {
		defer close(exitChan)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					defer attachResp.Close()
				}
				break
			}
			attachResp.Conn.Write(message)
		}
	}()

	<-exitChan
	log.Println("[Docker] terminal session close")

	return nil
}
