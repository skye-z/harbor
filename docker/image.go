package docker

import "github.com/docker/docker/api/types"

// 获取镜像列表
func (d Docker) GetImageList() ([]types.ImageSummary, error) {
	images, err := d.Session.ImageList(d.Context, types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return images, nil
}
