/*
Docker 日志解析工具

BetaX Harbor
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package util

import (
	"bufio"
	"io"
	"log"
	"strings"
)

func ProcessLogs(logs io.Reader) []string {
	var cleanedLogs []string
	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		line := extractLogContent(scanner.Text())
		cleanedLogs = append(cleanedLogs, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return cleanedLogs
}

func extractLogContent(s string) string {
	parts := strings.SplitN(s, " ", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return s
}
