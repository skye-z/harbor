package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type VersionInfo struct {
	Id      int64        `json:"id"`
	Url     string       `json:"html_url"`
	Title   string       `json:"name"`
	Content string       `json:"body"`
	Version string       `json:"tag_name"`
	Date    string       `json:"published_at"`
	Package []AppPackage `json:"assets"`
}

type AppPackage struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Url  string `json:"browser_download_url"`
}

func GetLatestReleaseVersion() *VersionInfo {
	resp, err := http.Get("https://api.github.com/repos/skye-z/harbor/releases/latest")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var info VersionInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil
	}

	return &info
}

func DownloadNewVersion() bool {
	info := GetLatestReleaseVersion()
	list := info.Package
	name := fmt.Sprintf("harbor_%s_%s", runtime.GOOS, runtime.GOARCH)
	url := ""
	for _, value := range list {
		if value.Name == name {
			url = value.Url
			break
		}
	}
	if url == "" {
		return false
	}

	log.Println("[Update] dwnload update file")
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	out, err := os.Create("update.cache")
	if err != nil {
		return false
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return false
	}
	log.Println("[Update] apply update")
	// 重命名文件
	if err := os.Rename("update.cache", "harbor"); err != nil {
		return false
	}
	return true
}

func RestartWithSystemd() error {
	log.Println("[Core] reboot")
	cmd := exec.Command("systemctl", "restart", "harbor")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
