package metrics_collector

type MetricsCollector interface {
	// StartPrometheusMetricsServer 启动prometheus抓取的服务
	StartPrometheusMetricsServer() error
}
