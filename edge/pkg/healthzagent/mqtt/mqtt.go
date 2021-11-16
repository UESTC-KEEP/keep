package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wonderivan/logger"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"time"
)

// Message  用于解析未知的mqtt数据
var Message map[string]interface{}
var MqttClient *client.Client

// SubscribeMqtt  用户传入需要监听的topic持续获取数据
func SubscribeMqtt(host_ip, port, topic string) map[string]interface{} {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	defer close(done)
	// build actual signal list to control

	connectToMqtt(host_ip, port)
	err := MqttClient.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte(topic),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					fmt.Println(string(topicName), string(message))
					err := json.Unmarshal(message, &Message)
					if err != nil {
						logger.Error(err)
					}
					fmt.Println("解析的结构体：", Message)
					cancel()
				},
			},
		},
	})
	if err != nil && ctx.Err() == nil {
		logger.Error(err)
	}

	select {
	// Check for termination request.
	case <-ctx.Done():

		fmt.Println("----------------------------")
		logger.Debug(fmt.Sprintf("Termination pending: %s", ctx.Err()))
		break
	case <-time.After(10000 * time.Millisecond):
	}

	//err := cli.Disconnect()
	if err != nil {
		logger.Error(err)
	}
	defer cancel()
	//defer MqttClient.Terminate()
	return Message
}

// GetRencentMqttMsg 获取当前最新的该mqtt主机的信息
func GetRencentMqttMsg(host_ip, port, topic string) map[string]interface{} {
	return SubscribeMqtt(host_ip, port, topic)
}

func connectToMqtt(host_ip, port string) {
	if MqttClient == nil {
		MqttClient = client.New(&client.Options{
			ErrorHandler: func(err error) {
				logger.Error("连接mqttbroker失败...", err)
			},
		})
		defer MqttClient.Terminate()
		// Connect to the MQTT Server.
		err := MqttClient.Connect(&client.ConnectOptions{
			Network:  "tcp",
			Address:  host_ip + ":" + port,
			ClientID: []byte("receive-client"),
		})
		if err != nil {
			logger.Fatal(err)
			panic(err)
		}
	}

}
