package metrics_server

type MetricServer interface {
	// CheckCadvisorStatus 检查集群中cadvisor的情况
	/*
		传入参数：汲取中master的ip:port字符串数组
	*/
	CheckCadvisorStatus(masters []string) error
	// StartCadvisorPod 启动cadvisor的组件 编写yaml
	/*
		传入参数：cadvisor的配置文件
	*/
	StartCadvisorPod(yamlFile string) error
}
