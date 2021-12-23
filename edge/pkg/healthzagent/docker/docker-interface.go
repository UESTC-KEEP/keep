package docker

type ContainerMetrics struct {
	Name             string
	ID               string
	Cpu              float64
	Memory           float64
	MemoryPercentage float64
	Disk             map[string]float64
	Net              map[string]float64
}

type DockerInterface interface {
	// ListContainerMetrics 查询对应的所有container的使用资源情况
	ListContainerMetrics() (*map[string]ContainerMetrics, error)
}
