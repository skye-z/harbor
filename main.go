package main

import (
	"harbor/docker"
)

func main() {
	client, err := docker.NewDocker()
	if err != nil {
		panic(err)
	}

	// 获取容器列表
	// containers, err := client.GetContainerList()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("container list:")
	// for i := 0; i < len(containers); i++ {
	// 	fmt.Println(containers[i].Names[0])
	// }

	// 获取镜像列表
	// images, err := client.GetImageList()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("image list:")
	// for i := 0; i < len(images); i++ {
	// 	fmt.Println(images[i].RepoTags)
	// }
	defer client.Close()
}
