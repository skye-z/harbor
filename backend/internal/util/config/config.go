package config

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/spf13/viper"
)

const Version = "0.1.0"

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			createDefault()
		} else {
			panic("Config read failed: " + err.Error())
		}
	}
}

func Reload() {
	viper.ReadInConfig()
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
