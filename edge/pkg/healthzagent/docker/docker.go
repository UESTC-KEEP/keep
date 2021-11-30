package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"keep/pkg/util/loggerv1.0.0"
)

var clit *client.Client

// 初始化客户端
func init() {
	clit_, err := client.NewClientWithOpts()
	if err != nil {
		logger.Warn("初始化docker客户端失败...")
		return
	}
	clit = clit_
}

// GetAllDockerImages 查看节点镜像列表
func GetAllDockerImages() (*[]types.ImageSummary, error) {

	list, err := clit.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		logger.Warn(err)
		return nil, err
	}
	return &list, nil
}
