/*
镜像数据模型

BetaX Harbor
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package model

import (
	"time"

	"xorm.io/xorm"
)

type Image struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Context   string `json:"context"`
	Timestamp int64  `json:"timestamp"`
}

type ImageModel struct {
	DB *xorm.Engine
}

// 添加镜像构建
func (model ImageModel) AddImage(name string, context string) bool {
	image := &Image{
		Name:    name,
		Context: context,
	}
	currentTime := time.Now()
	image.Timestamp = currentTime.UnixMilli()
	_, err := model.DB.Insert(image)
	return err == nil
}

// 获取镜像构建列表
func (model ImageModel) GetImageList() ([]Image, error) {
	var list []Image
	err := model.DB.Cols("id", "name", "timestamp").Desc("id").Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 获取镜像构建
func (model ImageModel) GetImageInfo(id int64) (*Image, error) {
	image := &Image{
		Id: id,
	}
	has, err := model.DB.Get(image)
	if !has {
		return nil, err
	}
	return image, nil
}

// 编辑镜像构建
func (model ImageModel) EditImage(image *Image) bool {
	if image.Id == 0 {
		return false
	}
	_, err := model.DB.ID(image.Id).Update(image)
	return err == nil
}

// 删除镜像构建
func (model ImageModel) DelImage(id int64) bool {
	if id == 0 {
		return false
	}
	_, err := model.DB.Delete(&Image{
		Id: id,
	})
	return err == nil
}
