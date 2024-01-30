package docker

import (
	"bytes"
	"fmt"
	"harbor/util"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/gin-gonic/gin"
)

type ImageBuild struct {
	Id       int    `json:"id"`
	Tag      string `json:"tag"`
	Store    string `json:"store"`
	Platform string `json:"platform"`
	Context  string `json:"context"`
}

// 获取镜像列表
func (d Docker) GetImageList() ([]types.ImageSummary, error) {
	images, err := d.Session.ImageList(d.Context, types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return images, nil
}

// 删除镜像
func (d Docker) RemoveImage(id string, prune bool, force bool) error {
	_, err := d.Session.ImageRemove(d.Context, id, types.ImageRemoveOptions{
		PruneChildren: prune,
		Force:         force,
	})
	return err
}

// 拉取镜像
func (d Docker) PullImage(ctx *gin.Context, form ImageBuild) {
	log.Printf("[Image] pull %s(%s) for %s\n", form.Tag, form.Platform, form.Store)
	out, err := d.Session.ImagePull(d.Context, form.Store+"/"+form.Tag, types.ImagePullOptions{
		Platform: form.Platform,
	})
	if out == nil || err != nil {
		util.ReturnMessageData(ctx, false, "镜像拉取失败", err.Error())
		log.Println("[Image] pull failed:", err)
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			log.Println("[Image] close failed:", closeErr)
		}
	}()
	ctx.Status(http.StatusOK)
	_, err = io.Copy(ctx.Writer, out)
	if err != nil {
		util.ReturnMessageData(ctx, false, "镜像拉取失败", err.Error())
		log.Printf("[Image] pull %s pulled failed\n", form.Tag)
	} else {
		util.ReturnMessage(ctx, true, "镜像拉取成功")
		log.Printf("[Image] pull %s pulled successfully\n", form.Tag)
	}
}

// 打标签
func (d Docker) AddImageTag(id string, tag string) error {
	err := d.Session.ImageTag(d.Context, id, tag)
	return err
}

// 搜索镜像
func (d Docker) SearchImage(keyword string) ([]registry.SearchResult, error) {
	list, err := d.Session.ImageSearch(d.Context, keyword, types.ImageSearchOptions{
		Limit: 20,
	})
	if err != nil {
		return []registry.SearchResult{}, err
	}
	return list, nil
}

// 获取镜像详情
func (d Docker) GetImageInfo(id string) (types.ImageInspect, error) {
	info, _, err := d.Session.ImageInspectWithRaw(d.Context, id)
	if err != nil {
		return types.ImageInspect{}, err
	}

	return info, nil
}

// 获取镜像构建历史
func (d Docker) GetImageHistory(id string) ([]image.HistoryResponseItem, error) {
	list, err := d.Session.ImageHistory(d.Context, id)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// 导出镜像
func (d Docker) ExportImage(ctx *gin.Context, id string) error {
	imageSaveReader, err := d.Session.ImageSave(d.Context, []string{id})
	if err != nil {
		util.ReturnMessage(ctx, false, "镜像导出失败")
		return err
	}
	defer imageSaveReader.Close()
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.tar", id))
	ctx.Header("Content-Type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, imageSaveReader)
	if err != nil {
		util.ReturnMessage(ctx, false, "镜像导出传输失败")
		return err
	}
	return nil
}

// 构建镜像
func (d Docker) BuildImage(tag, context string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(currentDir, "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	dockerfilePath := filepath.Join(cacheDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(context), 0644); err != nil {
		return "", err
	}

	build, err := d.Session.ImageBuild(d.Context, bytes.NewReader(nil), types.ImageBuildOptions{
		Dockerfile: dockerfilePath,
		Tags:       []string{tag},
	})
	if err != nil {
		return "", err
	}
	defer build.Body.Close()

	buildResponse, err := io.ReadAll(build.Body)
	if err != nil {
		return "", err
	}
	return string(buildResponse), nil
}
