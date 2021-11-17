package mqtt

import (
	"github.com/wonderivan/logger"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

const MQTT_QOS = mqtt.QoS0
const MQTT_CHAN_SIZE = 4
const LOG_TAG = "MQTT"

type MqttClient_t struct {
	p_mqtt_client *client.Client
	topic_map     map[string]chan []byte
}

func CreateMqttClient(client_name string, broker_ip string, broker_port string) *MqttClient_t {

	var mqtt_cli = client.New(&client.Options{
		ErrorHandler: func(err error) {
			logger.Error("连接mqttbroker失败...", err)
		},
	})
	defer mqtt_cli.Terminate()

	conn_opt := client.ConnectOptions{
		Network:  "tcp",
		Address:  broker_ip + ":" + broker_port,
		ClientID: []byte(client_name),
	}
	err := mqtt_cli.Connect(&conn_opt)

	if nil != err {
		logger.Fatal(err)
		panic(err)
	}

	p_cli := new(MqttClient_t)

	p_cli.p_mqtt_client = mqtt_cli
	p_cli.topic_map = make(map[string]chan []byte)

	return p_cli
}

//只要调用一次就行，其后等着自己的回调函数就行，不用反复注册订阅
func (mqtt_cli *MqttClient_t) RegistSubscribeTopic(topic string) {
	p_cli := mqtt_cli.p_mqtt_client

	_, exist := mqtt_cli.topic_map[topic]
	if exist { //不能重复订阅同一主题
		logger.Warn(LOG_TAG + ": Skip subscribeing duplicated topic " + topic)
		return
	}

	err := p_cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			{
				TopicFilter: []byte(topic),
				QoS:         MQTT_QOS,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					logger.Debug(LOG_TAG+": Topic= "+string(topicName)+"\tData=", message)
					mqtt_cli.topic_map[string(topicName)] <- message //TODO 这个地方得先设法获取通道是否满了，不然会在满了后一直阻塞
				},
			},
		},
	})
	if nil != err {
		logger.Error("Failed to subscribe ", err)
	}
	mqtt_cli.topic_map[topic] = make(chan []byte, MQTT_CHAN_SIZE)
}

func (mqtt_cli *MqttClient_t) GetTopicData(topic string) []byte {
	val, exist := mqtt_cli.topic_map[topic]
	if exist {
		data := <-val
		return data
	} else {
		logger.Error(LOG_TAG + "try to read unregisted topic " + topic)

		return nil
	}
}

func (mqtt_cli *MqttClient_t) UnSubscribeTopic(topic string) { //TODO 没做完
	logger.Fatal("unimplemented function")
	panic(nil)
}

func (mqtt_cli *MqttClient_t) PublishMsg(topic string, data []byte) {
	logger.Fatal("unimplemented function")
	panic(nil)
}
