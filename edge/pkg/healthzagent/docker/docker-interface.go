package docker

type ContainerMetrics struct {
	Cpu    float64
	Memory float64
	Disk   map[string]float64
	Net    map[string]int
}

type DockerInterface interface {
	// ListContainerMetrics 查询对应的所有container的使用资源情况
	ListContainerMetrics() (*map[string]ContainerMetrics, error)
}
