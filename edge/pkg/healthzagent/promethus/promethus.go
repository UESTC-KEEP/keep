package prome

import (
	"encoding/json"
	"fmt"
	"keep/constants"
	"keep/edge/pkg/healthzagent/config"
	"keep/edge/pkg/healthzagent/mqtt"
	"keep/edge/pkg/healthzagent/server"
	"keep/pkg/util/kelogger"
	"net/http"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

const LOG_TAG = "PROMETHUS"

func UnmarshalMqttData(data []byte) string {
	type JSONData_t map[string]string
	var msg JSONData_t
	err := json.Unmarshal(data, &msg)
	if err != nil {
		kelogger.Error(err)
	}
	kelogger.Debug(LOG_TAG+": 解析的结构体：", msg)
	strTemp := msg["temp"]
	return strTemp
}

var mqttCli *mqtt.MqttClient

func StartMertricsServer(port int) {
	//temp.
	// 暴露指标
	clientName := (uuid.NewV4()).String()
	mqttCli = mqtt.CreateMqttClient(clientName, constants.DefaultTestingMQTTServer, strconv.Itoa(constants.DefaultTestingMQTTPort))

	http.HandleFunc("/metrics", GetMetricOfEdge)
	kelogger.Debug(LOG_TAG + ": metricsServer启动成功...")
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		kelogger.Error(err)
	}
}

type Mtrics struct {
	Devices map[string]interface{} `json:"devices"`
	Metric  map[string]string      `json:"metric"`
}

func GetMetricOfEdge(resp http.ResponseWriter, req *http.Request) {
	kelogger.Debug("请求方：", req.RemoteAddr)
	DeviceMqttTopic := config.Config.DeviceMqttTopics

	retMap := Mtrics{}
	retMap.Metric = make(map[string]string)
	retMap.Devices = make(map[string]interface{})
	hzat := server.Healagent
	for i := 0; i < len(DeviceMqttTopic); i++ {
		fmt.Println(DeviceMqttTopic[i])
		mqttCli.RegistSubscribeTopic(&mqtt.TopicConf{TopicName: DeviceMqttTopic[i], TimeoutMs: 5000})
		// 设置 gauge 的值为
		dataRec, err := mqttCli.GetTopicData(DeviceMqttTopic[i]) //直接获取二进制数据，GetTopicData本身不做解析
		// mqttCli.UnSubscribeTopic(DeviceMqttTopic[i])
		if nil != err {
			kelogger.Error(LOG_TAG + ": Read mqtt err")
			time.Sleep(5 * time.Second)
		}
		tempData := UnmarshalMqttData(dataRec)
		newTemp, err := strconv.ParseFloat(tempData, 64)
		if err != nil {
			kelogger.Error(err)
			return
		}
		tempFloat, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", newTemp), 64)
		retMap.Devices[DeviceMqttTopic[i]] = tempFloat
	}
	retMap.Metric = map[string]string{"cpu": fmt.Sprintf("%.2f", hzat.CpuUsage), "mem": fmt.Sprintf("%.2f", hzat.Mem.UsedPercent)}
	resp.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(resp).Encode(&retMap)
	if err != nil {
		kelogger.Error(err)
	}
}
