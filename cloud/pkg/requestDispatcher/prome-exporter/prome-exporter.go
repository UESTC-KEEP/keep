package prome_exporter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_cpu_usage",
		Help: "Current usage of the CPU.",
	})
	memUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_mem_usage",
		Help: "Current usage of the Mem.",
	})

	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)
)
var newestCpuUsage = new(float64)
var newestMemUsage = new(float64)
var NodeExporterContentStr = ""

func StartPromeExporter() {
	cpuUsage.Set(*newestCpuUsage)
	hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/node-exporter", nodeexporter)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func nodeexporter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, NodeExporterContentStr)
}

// UpdatesData 更新数据为节点上传的最新数据
//func UpdatesData(nodeName string, cpu_usage float64, mem_usage float64) {
//	logger.Trace("更新节点用量数据:", nodeName, " cpu:", cpu_usage, " 内存：", mem_usage)
//	*newestCpuUsage = cpu_usage
//	*newestMemUsage = mem_usage
//	cpuUsage.Set(*newestCpuUsage)
//	memUsage.Set(*newestMemUsage)
//}

func UpdateNodeExporter(node_exporter string) {
	NodeExporterContentStr = node_exporter
}

func init() {
	*newestCpuUsage = 0.0
	*newestMemUsage = 0.0
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memUsage)
	prometheus.MustRegister(hdFailures)
}
