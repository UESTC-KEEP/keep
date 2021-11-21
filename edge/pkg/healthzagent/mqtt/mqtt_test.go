package mqtt

import (
	"testing"
	"time"

	"github.com/wonderivan/logger"
)

func TestSubscribeMqtt(t *testing.T) {
	tempTopic := "clock_sensor"
	mqttCli := CreateMqttClientNoName("192.168.1.40", "1883")
	mqttCli.RegistSubscribeTopic(&TopicConf{TopicName: "clock_sensor", TimeoutMs: 5000})
	for {
		dataRec, err := mqttCli.GetTopicData(tempTopic) //直接获取二进制数据，GetTopicData本身不做解析
		if nil != err {
			time.Sleep(time.Second) //如果试图读取不存在的主题，就会高速刷日志，在vscode下会大量占用内存，而且不接收信号
		} else {
			logger.Debug("TEST: mqtt rec=", string(dataRec))
		}
	}
}
