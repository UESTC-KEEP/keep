package prome

import (
	"encoding/json"
	"fmt"
	"keep/edge/pkg/healthzagent/mqtt"
	"keep/edge/pkg/healthzagent/server"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	uuid "github.com/satori/go.uuid"
	"github.com/wonderivan/logger"
)

const LOG_TAG = "PROMETHUS"

func UnmarshalMqttData(data []byte) string {
	type JSONData_t map[string]string
	var msg JSONData_t
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(LOG_TAG+": 解析的结构体：", msg)
	str_temp := string(msg["temp"])
	return str_temp
}

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
		temp_topic := "clock_sensor"
		client_name := (uuid.NewV4()).String()
		mqtt_cli := mqtt.CreateMqttClient(client_name, "192.168.1.40", "1883")
		mqtt_cli.RegistSubscribeTopic(&mqtt.MqttTopicConf{Topic_name: temp_topic, Timeout_ms: 5000})
		for {
			// 设置 gauge 的值为
			data_rec, err := mqtt_cli.GetTopicData(temp_topic) //直接获取二进制数据，GetTopicData本身不做解析
			if nil != err {                                    //TODO这个地方得考虑超时处理，算是检验设备是否在线的一部分
				logger.Error(LOG_TAG + ": Read mqtt err")
				time.Sleep(5 * time.Second)
				continue
			}

			temp_data := UnmarshalMqttData(data_rec)

			newTemp, err := strconv.ParseFloat(temp_data, 64)

			if err != nil {
				logger.Error(err)
				continue
			}
			hzat := server.Healagent
			mem.Set(hzat.Mem.UsedPercent)
			cpu.Set(hzat.CpuUsage)
			temp.Set(newTemp)
			fmt.Println("promethus更新:温度 " + fmt.Sprintf("%.2f", newTemp) + "  cpu:" + fmt.Sprintf("%.2f", hzat.CpuUsage) + "  mem:" + fmt.Sprintf("%.2f", hzat.Mem.UsedPercent))
		}
	}()
	//temp.
	// 暴露指标
	http.Handle("/metrics", promhttp.Handler())
	logger.Debug(LOG_TAG + ": metricsServer启动成功...")
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		logger.Error(err)
	}
}
