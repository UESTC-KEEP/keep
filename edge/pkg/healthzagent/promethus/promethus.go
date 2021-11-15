package prome

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wonderivan/logger"
	"keep/edge/pkg/healthzagent/mqtt"
	"keep/edge/pkg/healthzagent/server"
	"net/http"
	"strconv"
	"time"
)

func StartMertricsServer(port int) {
	// 创建一个没有任何 label 标签的 gauge 指标
	temp := AssembleModel("keep_device_temp", "温度传感器的温度")
	mem := AssembleModel("keep_node_mem_percent", "node节点内存使用量")
	cpu := AssembleModel("keep_node_cpu_percent", "node节点cpu使用量")
	// 在默认的注册表中注册该指标
	prometheus.MustRegister(temp)
	prometheus.MustRegister(mem)
	prometheus.MustRegister(cpu)

	go func() {
		for {
			// 设置 gauge 的值为
			temp_ := mqtt.GetRencentMqttMsg("192.168.1.40", "1883", "clock_sensor")["temp"].(string)
			newTemp, err := strconv.ParseFloat(temp_, 64)
			if err != nil {
				logger.Error(err)
				continue
			}
			hzat := server.Healagent
			mem.Set(hzat.Mem.UsedPercent)
			cpu.Set(hzat.CpuUsage)
			temp.Set(newTemp)
			fmt.Println("promethus更新:温度 " + fmt.Sprintf("%.2f", newTemp) + "  cpu:" + fmt.Sprintf("%.2f", hzat.CpuUsage) + "  mem:" + fmt.Sprintf("%.2f", hzat.Mem.UsedPercent))
			time.Sleep(5 * time.Second)
		}
	}()
	//temp.
	// 暴露指标
	http.Handle("/metrics", promhttp.Handler())
	logger.Debug("metricsServer启动成功...")
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		logger.Error(err)
	}
}
