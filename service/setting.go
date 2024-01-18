package service

import (
	"harbor/util"
	"strings"

	"github.com/gin-gonic/gin"
)

type SettingService struct {
}

func NewSettingService() *SettingService {
	ss := new(SettingService)
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
	clientSecret := util.GetString("oauth2.clientSecret")
	if clientSecret != "" {
		clientSecret = "**********"
	}
	config := &SettingOAuth2{
		Enable:       util.GetBool("oauth2.enable"),
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
	var form SettingOAuth2
	if err := ctx.ShouldBindJSON(&form); err != nil {
		util.ReturnMessage(ctx, false, "传入数据无效")
		return
	}
	if form.Enable {
		util.Set("oauth2.enable", 1)
	} else {
		util.Set("oauth2.enable", 0)
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
