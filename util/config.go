/*
全局配置服务

BetaX Quest
Copyright © 2023 SkyeZhang <skai-zhang@hotmail.com>
*/

package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const Version = "1.0.1"

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			createDefault()
		} else {
			// 配置文件被找到，但产生了另外的错误
			fmt.Println(err)
		}
	}
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
	viper.WriteConfig()
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func CheckMapExist(key, value string) bool {
	item := viper.GetStringMap(key)
	return item != nil && item[value] != nil
}

func AddMap(key, value string, date int64) {
	item := viper.GetStringMap(key)
	if item == nil {
		item = make(map[string]interface{})
	}
	item[value] = date
	Set(key, item)
}

func DelMap(key, value string) {
	item := viper.GetStringMap(key)
	if item != nil {
		delete(item, value)
		Set(key, item)
	}
}

func createDefault() {
	log.Println("[Config] init default config")
	// 安装状态
	viper.SetDefault("basic.install", "0")
	// 安全防护
	viper.SetDefault("secure.qps", "10")
	viper.SetDefault("secure.blacklist", "")
	// OAuth2
	viper.SetDefault("oauth2.enable", "0")
	viper.SetDefault("oauth2.clientId", "")
	viper.SetDefault("oauth2.clientSecret", "")
	viper.SetDefault("oauth2.redirectUrl", "")
	viper.SetDefault("oauth2.authUrl", "")
	viper.SetDefault("oauth2.tokenUrl", "")
	viper.SetDefault("oauth2.userUrl", "")
	viper.SetDefault("oauth2.scopes", "")
	viper.SetDefault("oauth2.userIdKey", "")
	viper.SetDefault("oauth2.userNameKey", "")
	// 监控告警
	viper.SetDefault("alarm.enable", "0")
	viper.SetDefault("alarm.path", "")
	viper.SetDefault("alarm.interval", "3")
	viper.SetDefault("alarm.event", "")
	viper.SetDefault("alarm.loadThreshold", "0.9")
	viper.SetDefault("alarm.memoryThreshold", "90")
	viper.SetDefault("alarm.diskThreshold", "90")
	// 令牌密钥
	secret, err := generateSecret()
	if err != nil {
		panic(err)
	}
	viper.SetDefault("token.secret", secret)
	// 令牌有效期/小时
	viper.SetDefault("token.exp", 24)
	viper.SafeWriteConfig()
}

func generateSecret() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
