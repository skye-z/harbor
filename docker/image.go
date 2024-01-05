package docker

import (
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
)

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
func (d Docker) PullImage(ctx *gin.Context, id string, store string, name string, platform string) {
	log.Printf("[Image] pull %s(%s) for %s\n", name, platform, store)
	out, err := d.Session.ImagePull(d.Context, store+"/"+name, types.ImagePullOptions{
		Platform: platform,
	})
	if err != nil {
		log.Println("[Image] pull failed:", err)
		return
	}
	defer out.Close()
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Status(http.StatusOK)
	_, err = io.Copy(ctx.Writer, out)
	if err != nil {
		log.Printf("[Image] pull %s pulled failed\n", name)
	} else {
		log.Printf("[Image] pull %s pulled successfully\n", name)
	}
}

// 打标签
func (d Docker) AddImageTag(id string, tag string) error {
	err := d.Session.ImageTag(d.Context, id, tag)
	return err
}

// 获取镜像详情
func (d Docker) GetImageInfo(id string) (types.ImageInspect, error) {
	info, _, err := d.Session.ImageInspectWithRaw(d.Context, id)
	if err != nil {
		return types.ImageInspect{}, err
	}

	return info, nil
}
