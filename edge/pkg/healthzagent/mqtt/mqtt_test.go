package mqtt

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/wonderivan/logger"
)

func TestSubscribeMqtt(t *testing.T) {
	temp_topic := "clock_sensor"
	client_name := (uuid.NewV4()).String()
	mqtt_cli := CreateMqttClient(client_name, "192.168.1.40", "1883")
	mqtt_cli.RegistSubscribeTopic("clock_sensor")
	for {
		// 设置 gauge 的值为
		data_rec := mqtt_cli.GetTopicData(temp_topic) //直接获取二进制数据，GetTopicData本身不做解析
		if nil == data_rec {
			time.Sleep(time.Second) //如果试图读取不存在的主题，就会高速刷日志，在vscode下会大量占用内存，而且不接收信号
		} else {
			logger.Debug("TEST: mqtt rec=", string(data_rec))
		}
	}
}
