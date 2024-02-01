package util

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

type BuildLogConfig struct {
	Type   string            `json:"Type"`
	Config map[string]string `json:"Config"`
}

type BuildRestartPolicy struct {
	Name              string `json:"Name"`
	MaximumRetryCount int    `json:"MaximumRetryCount"`
}

type BuildResources struct {
	CPU    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

type BuildContainer struct {
	Name      string                   `json:"name"`
	Container container.Config         `json:"container"`
	Host      HostConfig               `json:"host"`
	Network   network.EndpointSettings `json:"network"`
}

type HostConfig struct {
	Binds         []string           `json:"binds"`         // 容器与主机之间的卷绑定列表
	LogConfig     BuildLogConfig     `json:"logConfig"`     // 容器日志配置
	NetworkMode   string             `json:"networkMode"`   // 容器使用的网络模式
	PortBindings  nat.PortMap        `json:"portBindings"`  // 容器内部端口与主机端口的映射
	RestartPolicy BuildRestartPolicy `json:"restartPolicy"` // 容器的重启策略
	AutoRemove    bool               `json:"autoRemove"`    // 容器退出时是否自动移除
	Annotations   map[string]string  `json:"annotations"`   // 附加到容器的任意非标识元数据

	// UNIX platforms
	CapAdd     strslice.StrSlice `json:"capAdd"`     // 添加到容器的内核功能列表
	CapDrop    strslice.StrSlice `json:"capDrop"`    // 从容器中移除的内核功能列表
	Privileged bool              `json:"privileged"` // 容器是否运行在特权模式下
	Resources  BuildResources    `json:"resources"`  // 容器资源配置

	// Mounts
	Mounts []mount.Mount `json:"mounts"` // 容器挂载配置
}
