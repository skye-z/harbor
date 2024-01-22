package monitor

import (
	"harbor/util"
	"net/http"
	"net/url"
	"strings"
)

func SendNotice(msg string) bool {
	if !util.GetBool("alarm.enable") {
		return false
	}
	url := util.GetString("alarm.path") + msg
	response, err := http.Get(url)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode == http.StatusOK
}

func SendEventNotice(msg string, event string) bool {
	events := strings.Split(util.GetString("alarm.event"), ",")
	found := false
	for _, word := range events {
		if word == event {
			found = true
			break
		}
	}
	if !found {
		return false
	}
	return SendNotice(url.QueryEscape(msg))
}

func NoticeContainerStart(id string) bool {
	return SendEventNotice("Harbor - 容器启动\t#"+id, "containerStart")
}

func NoticeContainerStop(id string) bool {
	return SendEventNotice("Harbor - 容器停止\t#"+id, "containerStop")
}

func NoticeDaemonShutdown() bool {
	return SendEventNotice("Harbor - 引擎下线\t警告!! Docker守护进程下线, Docker引擎停止服务!!!", "daemonShutdown")
}

func NoticeHighLoad() bool {
	return SendEventNotice("Harbor - 突发高负载\t警告! 检测到负载过高, 请及时进行处理", "highLoad")
}

func NoticeRunOut(name string, number string) bool {
	return SendEventNotice("Harbor - 资源即将耗尽\t警告! "+name+"已使用"+number+"%, 请及时进行处理", "runOut")
}
