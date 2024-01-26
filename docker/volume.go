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

// 获取存储详情
func (d Docker) GetVolumeInfo(id string) (volume.Volume, error) {
	info, err := d.Session.VolumeInspect(d.Context, id)
	if err != nil {
		return volume.Volume{}, err
	}

	return info, nil
}

// 创建存储
func (d Docker) CreateVolume(name string, driver string) error {
	_, err := d.Session.VolumeCreate(d.Context, volume.CreateOptions{
		Name:   name,
		Driver: driver,
	})
	if err != nil {
		return err
	}

	return nil
}

// 删除存储
func (d Docker) RemoveVolume(id string) error {
	err := d.Session.VolumeRemove(d.Context, id, false)
	return err
}
