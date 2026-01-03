package config

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/spf13/viper"
)

const Version = "2.0.0"

var ConfigPath = "."

func InitConfig() {
	InitConfigWithPath(".")
}

func InitConfigWithPath(path string) {
	ConfigPath = path
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")

	if path != "." {
		viper.AddConfigPath("/opt/harbor")
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			createDefault()
		} else {
			panic("Config read failed: " + err.Error())
		}
	}
}

func Reload() error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
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

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func createDefault() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("cors.origins", "")
	viper.SetDefault("docker.socket", "/var/run/docker.sock")

	secret, err := generateSecret()
	if err != nil {
		panic("Failed to generate secret: " + err.Error())
	}
	viper.SetDefault("jwt.secret", secret)
	viper.SetDefault("jwt.expiration", "24h")

	if ConfigPath == "/opt/harbor" || ConfigPath == "." {
		viper.SafeWriteConfigAs("/opt/harbor/config.yaml")
	} else {
		viper.SafeWriteConfig()
	}
}

func generateSecret() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
