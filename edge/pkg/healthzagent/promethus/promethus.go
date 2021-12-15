package prome

import (
	"encoding/json"
	"fmt"
	"keep/edge/pkg/healthzagent/mqtt"
	"keep/edge/pkg/healthzagent/server"
	"keep/pkg/util/kplogger"
	"net/http"
	"strconv"
)

const LOG_TAG = "<PROMETHUS>"

type Metrics_t struct {
	Devices map[string]interface{} `json:"devices"`
	Metric  map[string]string      `json:"metric"`
}

func UnmarshalMqttData(data []byte) string {
	type JSONData_t map[string]string
	var msg JSONData_t
	err := json.Unmarshal(data, &msg)
	if err != nil {
		kplogger.Error(err)
	}
	kplogger.Debug(LOG_TAG+": 解析的结构体：", msg)
	strTemp := msg["temp"]
	return strTemp
}

var mqttCli *mqtt.MqttClient = nil

//func InitMqttClient() {
//	mqttCli = mqtt.CreateMqttClientNoName(edge.DefaultTestingMQTTServer, strconv.Itoa(edge.DefaultTestingMQTTPort))
//	DeviceMqttTopic := config.Config.DeviceMqttTopics
//	for i := 0; i < len(DeviceMqttTopic); i++ {
//		//MQTT_CACHE_MODE不会阻塞当前协程，而是返回最新缓存的数据，不一定是当前时刻的
//		mqttCli.RegistSubscribeTopic(&mqtt.TopicConf{TopicName: DeviceMqttTopic[i], TimeoutMs: 5000, DataMode: mqtt.MQTT_CACHE_MODE})
//	}
//}

func StartMertricsServer(port int) {
	//temp.
	// 暴露指标
	//InitMqttClient()
	http.HandleFunc("/metrics", reportMetricOfEdge)

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		kplogger.Error(err)
	}
	kplogger.Debug(LOG_TAG + ": metricsServer启动成功...")
}

func StopMertricsServer() {
	if nil != mqttCli {

		mqttCli.DestroyMqttClient()
		//TODO  应该还有停止http的东西
	}
}

func reportMetricOfEdge(resp http.ResponseWriter, req *http.Request) {
	//kplogger.Debug("请求方：", req.RemoteAddr)
	//DeviceMqttTopic := config.Config.DeviceMqttTopics

	retMap := Metrics_t{}
	retMap.Metric = make(map[string]string)
	retMap.Devices = make(map[string]interface{})
	hzat := server.Healagent

	//for i := 0; i < mqttCli.GetTopicNum(); i++ { //TODO 直接取序号肯定不合适， 这里要做成map的迭代器，但是不能暴露map，
	//	fmt.Println(DeviceMqttTopic[i])
	//
	//	dataRec, err := mqttCli.GetTopicData(DeviceMqttTopic[i]) //直接获取二进制数据，GetTopicData本身不做解析
	//
	//	if nil != err {
	//		kplogger.Error(LOG_TAG+": Read mqtt err", err.Error())
	//		time.Sleep(time.Millisecond * 100) //TODO 时间有待调整，或者取消
	//		continue
	//	}
	//	tempData := UnmarshalMqttData(dataRec)
	//	newTemp, err := strconv.ParseFloat(tempData, 64)
	//	if err != nil {
	//		kplogger.Error(err)
	//		return
	//	}
	//	tempFloat, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", newTemp), 64)
	//	retMap.Devices[DeviceMqttTopic[i]] = tempFloat
	//}
	retMap.Metric = map[string]string{"cpu": fmt.Sprintf("%.2f", hzat.CpuUsage), "mem": fmt.Sprintf("%.2f", hzat.Mem.UsedPercent)}
	resp.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(resp).Encode(&retMap)
	if err != nil {
		kplogger.Error(err)
	}
}
