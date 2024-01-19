package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func CalculateMD5(input string) string {
	data := []byte(input)
	hash := md5.Sum(data)
	hashInHex := hex.EncodeToString(hash[:])
	return hashInHex
}
