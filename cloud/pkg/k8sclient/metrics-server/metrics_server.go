package metrics_server

import (
	"context"
	"fmt"
	"keep/cloud/pkg/k8sclient/config"
	"time"
)

type MetricServerImpl struct{}

func (msi *MetricServerImpl) CheckCadvisorStatus(masters []string) error {
	//var wg sync.WaitGroup
	for _, master := range masters {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.MasterMetricTimeout)*time.Millisecond)
		fmt.Println(time.Duration(config.Config.MasterMetricTimeout) * time.Millisecond)
		done := make(chan bool)
		go func(master string) {
			select {
			case <-done:
				fmt.Println("查询完成...")
				cancel()
				return
			case <-ctx.Done():
				fmt.Println(master + "超时......")
			}
		}(master)
		go getMasterMetrics(master, done)
	}
	for {
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (msi *MetricServerImpl) StartCadvisorPod(yamlFile string) error {

	return nil
}

func NewMetricServerImpl() *MetricServerImpl {
	return &MetricServerImpl{}
}

func getMasterMetrics(master string, done chan bool) {
	fmt.Println("开始查询master:" + master + " 数据...")
	time.Sleep(10 * time.Second)
	done <- true
}
