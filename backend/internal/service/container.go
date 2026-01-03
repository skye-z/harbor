package service

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/skye-z/harbor/internal/util/docker"
	"github.com/skye-z/harbor/internal/util/response"
)

// 容器管理服务
type ContainerService struct {
	client *docker.Client
}

// 创建容器服务实例
func NewContainerService(client *docker.Client) *ContainerService {
	return &ContainerService{client: client}
}

// 获取容器列表
func (s *ContainerService) GetList(c *gin.Context) {
	containers, err := s.client.ListContainers()
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, containers)
}

// 获取容器的详细信息
func (s *ContainerService) GetInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	info, err := s.client.GetContainerInfo(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, info)
}

// 容器操作常量
const (
	ContainerActionStart   = 1
	ContainerActionStop    = 2
	ContainerActionRestart = 3
	ContainerActionRemove  = 4
	ContainerActionPause   = 5
	ContainerActionUnpause = 6
)

// 容器操作
func (s *ContainerService) Operation(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	action := c.DefaultQuery("action", "start")

	var actionCode int
	switch action {
	case "start":
		actionCode = ContainerActionStart
	case "stop":
		actionCode = ContainerActionStop
	case "restart":
		actionCode = ContainerActionRestart
	case "remove":
		actionCode = ContainerActionRemove
	case "pause":
		actionCode = ContainerActionPause
	case "unpause":
		actionCode = ContainerActionUnpause
	default:
		response.BadRequest(c, "不支持的操作类型，支持: start/stop/restart/remove/pause/unpause")
		return
	}

	err := s.client.OperationContainer(id, actionCode)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "操作成功", gin.H{"action": action})
}

// 获取容器日志
func (s *ContainerService) GetLogs(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	tail := c.DefaultQuery("tail", "100")
	logs, err := s.client.GetContainerLogs(c.Request.Context(), id, tail)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	defer logs.Close()

	content, _ := io.ReadAll(logs)
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, string(content))
}

// 获取容器统计信息
func (s *ContainerService) GetStat(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	stats, err := s.client.GetContainerStats(id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// 获取容器进程信息
func (s *ContainerService) GetProcesses(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	processes, err := s.client.GetContainerProcesses(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, processes)
}

// 获取容器文件目录
func (s *ContainerService) GetDiff(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	path := c.DefaultQuery("path", "/")
	result, err := s.client.ListContainerDir(c.Request.Context(), id, path)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, result)
}

// 从容器复制文件
func (s *ContainerService) CopyFromContainer(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	srcPath := c.Query("src_path")
	if srcPath == "" {
		response.BadRequest(c, "缺少源路径")
		return
	}

	dstPath := c.DefaultQuery("dst_path", "./tmp")
	err := s.client.CopyFromContainer(c.Request.Context(), id, srcPath, dstPath)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "文件已复制到本地", nil)
}

// 复制文件到容器
func (s *ContainerService) CopyToContainer(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	srcPath := c.Query("src_path")
	if srcPath == "" {
		response.BadRequest(c, "缺少源路径")
		return
	}

	dstPath := c.Query("dst_path")
	if dstPath == "" {
		response.BadRequest(c, "缺少目标路径")
		return
	}

	if strings.Contains(dstPath, "..") || strings.HasPrefix(dstPath, "/") {
		response.BadRequest(c, "非法的目标路径")
		return
	}

	err := s.client.CopyToContainer(c.Request.Context(), id, srcPath, dstPath)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "文件已上传到容器", nil)
}

// 打开容器终端
func (s *ContainerService) ConnectTerminal(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	execID, err := s.client.CreateExec(c.Request.Context(), id, &docker.ContainerCreateConfig{
		Cmd:          []string{"/bin/sh"},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	})
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	c.Header("X-Exec-ID", execID)
	response.Success(c, gin.H{
		"exec_id": execID,
		"ws_url":  "/api/container/terminal/ws?exec_id=" + execID,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 终端WebSocket连接
func (s *ContainerService) TerminalWebSocket(c *gin.Context) {
	execID := c.Query("exec_id")
	if execID == "" {
		response.BadRequest(c, "缺少执行实例ID")
		return
	}

	result, err := s.client.AttachToExec(c.Request.Context(), execID, false)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	defer result.Conn.Close()

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := result.Reader.Read(buf)
			if err != nil {
				break
			}
			ws.WriteMessage(websocket.TextMessage, buf[:n])
		}
		ws.Close()
	}()

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			break
		}
		result.Conn.Write(data)
	}
}

// 调整终端大小
func (s *ContainerService) ResizeTerminal(c *gin.Context) {
	execID := c.Query("exec_id")
	if execID == "" {
		response.BadRequest(c, "缺少执行实例ID")
		return
	}

	rowsStr := c.DefaultQuery("rows", "24")
	colsStr := c.DefaultQuery("cols", "80")
	rows, err := strconv.Atoi(rowsStr)
	if err != nil || rows <= 0 {
		rows = 24
	}
	cols, err := strconv.Atoi(colsStr)
	if err != nil || cols <= 0 {
		cols = 80
	}

	err = s.client.ExecResize(c.Request.Context(), execID, rows, cols)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "终端大小已调整", nil)
}

// 关闭容器终端
func (s *ContainerService) CloseTerminal(c *gin.Context) {
	execID := c.Query("exec_id")
	if execID == "" {
		response.BadRequest(c, "缺少执行实例ID")
		return
	}

	err := s.client.CloseExec(c.Request.Context(), execID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "终端已关闭", nil)
}

// 将容器打包为镜像
func (s *ContainerService) CommitContainer(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	repo := c.DefaultQuery("repo", "")
	tag := c.DefaultQuery("tag", "latest")

	imageID, err := s.client.CommitContainer(id, repo, tag)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "镜像构建成功", gin.H{
		"id":   imageID,
		"repo": repo,
		"tag":  tag,
	})
}

// 创建容器
func (s *ContainerService) CreateContainer(c *gin.Context) {
	imageName := c.Query("image")
	if imageName == "" {
		response.BadRequest(c, "缺少镜像名称")
		return
	}

	containerName := c.DefaultQuery("name", "")
	cmdStr := c.DefaultQuery("cmd", "")
	var cmd []string
	if cmdStr != "" {
		cmd = strings.Split(cmdStr, " ")
	}

	envStr := c.DefaultQuery("env", "")
	var env []string
	if envStr != "" {
		env = strings.Split(envStr, ",")
	}

	container, err := s.client.CreateContainer(imageName, containerName, cmd, env)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "容器创建成功", container)
}

// 重命名容器
func (s *ContainerService) RenameContainer(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少容器ID")
		return
	}

	newName := c.Query("name")
	if newName == "" {
		response.BadRequest(c, "缺少新名称")
		return
	}

	err := s.client.RenameContainer(id, newName)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "容器重命名成功", gin.H{
		"id":      id,
		"newName": newName,
	})
}
