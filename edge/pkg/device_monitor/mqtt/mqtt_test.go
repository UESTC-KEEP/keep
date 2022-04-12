package mqtt

import (
	"fmt"
	"github.com/UESTC-KEEP/keep/pkg/util/kplogger"
	"testing"
	"time"
)

func TestSubscribeMqttBlockMode(t *testing.T) {
	tempTopic := "clock_sensor"
	mqttCli := CreateMqttClientNoName("192.168.1.40", "1883")
	mqttCli.RegistSubscribeTopic(&TopicConf{TopicName: "clock_sensor", TimeoutMs: 5000, DataMode: MQTT_BLOCK_MODE})
	for {
		dataRec, err := mqttCli.GetTopicData(tempTopic) //直接获取二进制数据，GetTopicData本身不做解析
		if nil != err {
			fmt.Println(err.Error())
			time.Sleep(time.Second) //如果试图读取不存在的主题，就会高速刷日志，在vscode下会大量占用内存，而且不接收信号
		} else {
			kplogger.Debug("TEST: mqtt rec=", string(dataRec))
		}
	}
}

func TestSubscribeMqttCacheMode(t *testing.T) {
	tempTopic := "clock_sensor"
	mqttCli := CreateMqttClientNoName("192.168.1.40", "1883")
	mqttCli.RegistSubscribeTopic(&TopicConf{TopicName: "clock_sensor", TimeoutMs: 5000, DataMode: MQTT_CACHE_MODE})
	for {
		dataRec, err := mqttCli.GetTopicData(tempTopic) //直接获取二进制数据，GetTopicData本身不做解析
		if nil != err {
			kplogger.Error(err.Error())
			time.Sleep(time.Millisecond * 400) //如果试图读取不存在的主题，就会高速刷日志，在vscode下会大量占用内存，而且不接收信号
		} else {
			kplogger.Debug("TEST: mqtt rec=", string(dataRec))
			time.Sleep(time.Millisecond * 400)
		}
	}
}
