package prome

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/wonderivan/logger"
	"keep/constants"
	"keep/edge/pkg/healthzagent/config"
	"keep/edge/pkg/healthzagent/mqtt"
	"keep/edge/pkg/healthzagent/server"
	"net/http"
	"strconv"
	"time"
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
	strTemp := msg["temp"]
	return strTemp
}

func StartMertricsServer(port int) {
	//temp.
	// 暴露指标
	http.HandleFunc("/metrics", GetMetricOfEdge)
	logger.Debug(LOG_TAG + ": metricsServer启动成功...")
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		logger.Error(err)
	}
}

type Mtrics struct {
	Devices map[string]interface{} `json:"devices"`
	Metric  map[string]string      `json:"metric"`
}

func GetMetricOfEdge(w http.ResponseWriter, r *http.Request) {
	logger.Debug("请求方：", r.RemoteAddr)
	DeviceMqttTopic := config.Config.DeviceMqttTopics
	clientName := (uuid.NewV4()).String()
	mqttCli := mqtt.CreateMqttClient(clientName, constants.DefaultTestingMQTTServer, strconv.Itoa(constants.DefaultTestingMQTTPort))
	retMap := Mtrics{}
	retMap.Metric = make(map[string]string)
	retMap.Devices = make(map[string]interface{})
	hzat := server.Healagent
	for i := 0; i < len(DeviceMqttTopic); i++ {
		fmt.Println(DeviceMqttTopic[i])
		mqttCli.RegistSubscribeTopic(&mqtt.TopicConf{TopicName: DeviceMqttTopic[i], TimeoutMs: 5000})
		// 设置 gauge 的值为
		dataRec, err := mqttCli.GetTopicData(DeviceMqttTopic[i]) //直接获取二进制数据，GetTopicData本身不做解析
		mqttCli.UnSubscribeTopic(DeviceMqttTopic[i])
		if nil != err { //TODO这个地方得考虑超时处理，算是检验设备是否在线的一部分
			logger.Error(LOG_TAG + ": Read mqtt err")
			time.Sleep(5 * time.Second)
		}
		tempData := UnmarshalMqttData(dataRec)
		newTemp, err := strconv.ParseFloat(tempData, 64)
		if err != nil {
			logger.Error(err)
			return
		}
		tempFloat, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", newTemp), 64)
		retMap.Devices[DeviceMqttTopic[i]] = tempFloat
	}
	retMap.Metric = map[string]string{"cpu": fmt.Sprintf("%.2f", hzat.CpuUsage), "mem": fmt.Sprintf("%.2f", hzat.Mem.UsedPercent)}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&retMap)
	if err != nil {
		logger.Error(err)
	}
}
