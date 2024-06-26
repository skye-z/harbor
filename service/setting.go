/*
设置服务

BetaX Harbor
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package service

import (
	"encoding/json"
	"os/exec"
	"strings"

	githubreleases "github.com/skye-z/github-releases"
	"github.com/skye-z/harbor/util"

	"github.com/gin-gonic/gin"
)

const (
	OAuth2Enable = "oauth2.enable"
	AlarmEnable  = "alarm.enable"
)

type SettingService struct {
	Version *githubreleases.Versioning
}

func NewSettingService() *SettingService {
	ss := new(SettingService)
	ss.Version = &githubreleases.Versioning{
		Author: "skye-z",
		Store:  "harbor",
		Name:   "harbor",
		Cmd:    exec.Command("systemctl", "restart", "harbor"),
		Proxy:  "https://mirror.ghproxy.com/",
	}
	return ss
}

type SettingOAuth2 struct {
	Enable       bool   `json:"enable"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"redirectUrl"`
	AuthURL      string `json:"authUrl"`
	TokenURL     string `json:"tokenUrl"`
	UserURL      string `json:"userUrl"`
	UserIdKey    string `json:"userIdKey"`
	UserNameKey  string `json:"userNameKey"`
	Scopes       string `json:"scopes"`
}

func (ss SettingService) GetOAuth2Setting(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	clientSecret := util.GetString("oauth2.clientSecret")
	if clientSecret != "" {
		clientSecret = "**********"
	}
	config := &SettingOAuth2{
		Enable:       util.GetBool(OAuth2Enable),
		ClientID:     util.GetString("oauth2.clientId"),
		ClientSecret: clientSecret,
		RedirectURL:  util.GetString("oauth2.redirectUrl"),
		Scopes:       util.GetString("oauth2.scopes"),
		AuthURL:      util.GetString("oauth2.authUrl"),
		TokenURL:     util.GetString("oauth2.tokenUrl"),
		UserURL:      util.GetString("oauth2.userUrl"),
		UserIdKey:    util.GetString("oauth2.userIdKey"),
		UserNameKey:  util.GetString("oauth2.userNameKey"),
	}
	util.ReturnData(ctx, true, config)
}

func (ss SettingService) UpdateOAuth2Setting(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	var form SettingOAuth2
	if err := ctx.ShouldBindJSON(&form); err != nil {
		util.ReturnMessage(ctx, false, "传入数据无效")
		return
	}
	if form.Enable {
		util.Set(OAuth2Enable, 1)
	} else {
		util.Set(OAuth2Enable, 0)
	}
	util.Set("oauth2.clientId", form.ClientID)
	if !strings.HasPrefix(form.ClientSecret, "*") && !strings.HasSuffix(form.ClientSecret, "*") {
		util.Set("oauth2.clientSecret", form.ClientSecret)
	}
	util.Set("oauth2.redirectUrl", form.RedirectURL)
	util.Set("oauth2.scopes", form.Scopes)
	util.Set("oauth2.authUrl", form.AuthURL)
	util.Set("oauth2.tokenUrl", form.TokenURL)
	util.Set("oauth2.userUrl", form.UserURL)
	util.Set("oauth2.userIdKey", form.UserIdKey)
	util.Set("oauth2.userNameKey", form.UserNameKey)

	util.ReturnMessage(ctx, true, "更新成功")
}

type SettingAlarm struct {
	Enable          bool    `json:"enable"`
	Path            string  `json:"path"`
	Interval        int     `json:"interval"`
	Event           string  `json:"event"`
	LoadThreshold   float64 `json:"loadThreshold"`
	MemoryThreshold int     `json:"memoryThreshold"`
	DiskThreshold   int     `json:"diskThreshold"`
}

func (ss SettingService) GetAlarmSetting(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	config := &SettingAlarm{
		Enable:          util.GetBool(AlarmEnable),
		Path:            util.GetString("alarm.path"),
		Interval:        util.GetInt("alarm.interval"),
		Event:           util.GetString("alarm.event"),
		LoadThreshold:   util.GetFloat64("alarm.loadThreshold"),
		MemoryThreshold: util.GetInt("alarm.memoryThreshold"),
		DiskThreshold:   util.GetInt("alarm.diskThreshold"),
	}
	util.ReturnData(ctx, true, config)
}

func (ss SettingService) UpdateAlarmSetting(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	var form SettingAlarm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		util.ReturnMessage(ctx, false, "传入数据无效")
		return
	}
	if form.Enable {
		util.Set(AlarmEnable, 1)
	} else {
		util.Set(AlarmEnable, 0)
	}
	util.Set("alarm.path", form.Path)
	util.Set("alarm.interval", form.Interval)
	util.Set("alarm.event", form.Event)
	util.Set("alarm.loadThreshold", form.LoadThreshold)
	util.Set("alarm.memoryThreshold", form.MemoryThreshold)
	util.Set("alarm.diskThreshold", form.DiskThreshold)

	util.ReturnMessage(ctx, true, "更新成功")
}

type SettingSecure struct {
	Qps       string `json:"qps"`
	BlackList string `json:"blacklist"`
}

func (ss SettingService) GetSecureSetting(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	blacklist, err := json.Marshal(util.GetStringMap("secure.blacklist"))
	if err != nil {
		util.ReturnMessage(ctx, false, "数据错误")
		return
	}

	config := &SettingSecure{
		Qps:       util.GetString("secure.qps"),
		BlackList: string(blacklist),
	}
	util.ReturnData(ctx, true, config)
}

func (ss SettingService) UpdateSecureSetting(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	var form SettingSecure
	if err := ctx.ShouldBindJSON(&form); err != nil {
		util.ReturnMessage(ctx, false, "传入数据无效")
		return
	}

	var blacklist = make(map[string]interface{})
	err := json.Unmarshal([]byte(form.BlackList), &blacklist)
	if err != nil {
		util.ReturnMessage(ctx, false, "传入数据非法")
		return
	}

	util.Set("secure.qps", form.Qps)
	util.Set("secure.blacklist", blacklist)

	util.ReturnMessage(ctx, true, "更新成功")
}

func (ss SettingService) GetNewVersion(ctx *gin.Context) {
	info := ss.Version.GetLatestReleaseVersion()
	if info == nil {
		util.ReturnMessage(ctx, false, "获取版本信息失败")
	} else {
		util.ReturnData(ctx, true, info)
	}
}

func (ss SettingService) UpdateNewVersion(ctx *gin.Context) {
	state := ss.Version.DownloadNewVersion()
	if state {
		util.ReturnMessage(ctx, true, "更新成功")
		ss.Version.RestartWithSystemd()
	} else {
		util.ReturnMessage(ctx, false, "更新失败")
	}
}
