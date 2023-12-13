package docker

import "github.com/docker/docker/api/types/volume"

// 获取存储卷列表
func (d Docker) GetVolumeList() ([]*volume.Volume, error) {
	volumes, err := d.Session.VolumeList(d.Context, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	return volumes.Volumes, nil
}
