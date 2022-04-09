package metrics_collector

import (
	"github.com/UESTC-KEEP/keep/cloud/pkg/promserver/config"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"net/http"
	"strconv"
)

type MetricsCollectorImpl struct{}

func (mci *MetricsCollectorImpl) StartPrometheusMetricsServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", MetricAssembler)
	logger.Debug("promserver服务启动中...  :" + strconv.Itoa(config.Config.PromServerPrometheusPort) + " 服务启动中...")
	err := http.ListenAndServe(":"+strconv.Itoa(config.Config.PromServerPrometheusPort), mux)
	if err != nil {
		logger.Fatal("promserver 启动失败,err：", err)
	}
	return nil
}

func MetricAssembler(w http.ResponseWriter, r *http.Request) {
	// 从kafka获取的消息
	//var msg model.Message

}
