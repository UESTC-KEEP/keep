package mqtt

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/wonderivan/logger"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

//config
const MQTT_QOS = mqtt.QoS0
const MQTT_CHAN_SIZE = 4

const LOG_TAG = "<MQTT>"
const MQTT_FOREVER = 0 //0表示无限等待

type MqttErrCode_t int

const (
	MQTT_OK MqttErrCode_t = iota
	MQTT_CHAN_CLOSED
	MQTT_TIMEOUT
	MQTT_TOPIC_UNEXIST
)

type MqttErrRet_t struct {
	err_code MqttErrCode_t
}

func (err *MqttErrRet_t) Error() string {
	return LOG_TAG //TODO 以后写点详细信息
}

type MqttTopicInfo_t struct {
	data_chan     chan []byte
	time_limit_ms uint
}

type MqttTopicMap map[string]*MqttTopicInfo_t
type MqttClient_t struct {
	p_mqtt_client *client.Client
	topic_map     MqttTopicMap
}

type MqttTopicConf struct {
	Topic_name string
	Timeout_ms uint
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
	p_cli.topic_map = make(MqttTopicMap)

	return p_cli
}

func CreateMqttClientNoName(broker_ip string, broker_port string) *MqttClient_t { //随机生成客户端名字
	return CreateMqttClient((uuid.NewV4()).String(), broker_ip, broker_port)
}

//只要调用一次就行，其后等着自己的回调函数就行，不用反复注册订阅
func (mqtt_cli *MqttClient_t) RegistSubscribeTopic(p_conf *MqttTopicConf) {
	topic_name := p_conf.Topic_name
	p_cli := mqtt_cli.p_mqtt_client

	_, exist := mqtt_cli.topic_map[topic_name]
	if exist { //不能重复订阅同一主题
		logger.Warn(LOG_TAG + ": Skip subscribeing duplicated topic " + topic_name)
		return
	}

	err := p_cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			{
				TopicFilter: []byte(topic_name),
				QoS:         MQTT_QOS,
				Handler: func(topicName, message []byte) {
					logger.Trace(LOG_TAG+": Topic= "+string(topicName)+"\tData=", message)
					mqtt_cli.topic_map[string(topicName)].data_chan <- message //TODO 这个地方得先设法获取通道是否满了，不然会在满了后一直阻塞
				},
			},
		},
	})
	if nil != err {
		logger.Error(LOG_TAG, ":Failed to subscribe topic: ", topic_name, " error=", err)
		return
	}

	mqtt_cli.topic_map[topic_name] = &MqttTopicInfo_t{
		data_chan:     make(chan []byte, MQTT_CHAN_SIZE),
		time_limit_ms: p_conf.Timeout_ms,
	}
}

func (mqtt_cli *MqttClient_t) GetTopicData(topic string) ([]byte, error) {
	topic_info, exist := mqtt_cli.topic_map[topic]
	if exist {
		if MQTT_FOREVER == topic_info.time_limit_ms { //TODO 不知道怎么复用select，凑合一下
			data := <-topic_info.data_chan
			if data == nil {
				logger.Warn(LOG_TAG + ":The data channel of topic \"" + topic + "\" was closed")
				return nil, &MqttErrRet_t{MQTT_CHAN_CLOSED}
			} else {
				return data, nil
			}
		} else {

			select {
			case data := <-topic_info.data_chan:
				if data == nil {
					logger.Warn(LOG_TAG + ":The data channel of topic \"" + topic + "\" was closed")
					return nil, &MqttErrRet_t{MQTT_CHAN_CLOSED}
				} else {
					return data, nil
				}
			case <-time.After(time.Duration(topic_info.time_limit_ms) * time.Millisecond):
				logger.Error(LOG_TAG, ": TIMEOUT while reading topic: ", topic)
				return nil, &MqttErrRet_t{MQTT_TIMEOUT}
			}
		}

	} else {
		logger.Error(LOG_TAG + "try to read unregisted topic " + topic)
		return nil, &MqttErrRet_t{MQTT_TOPIC_UNEXIST}
	}
}

func (mqtt_cli *MqttClient_t) UnSubscribeTopic(topic string) {
	topic_info, exist := mqtt_cli.topic_map[topic]
	if exist {
		mqtt_cli.p_mqtt_client.Unsubscribe(
			&client.UnsubscribeOptions{
				TopicFilters: [][]byte{
					[]byte(topic)},
			})
		close(topic_info.data_chan)
		//chan 关闭时的原则是：不要在接收协程中关闭，并且，如果有多个发送者时就不要关闭chan了。
		//https://studygolang.com/articles/9478
		delete(mqtt_cli.topic_map, topic)
	} else {
		logger.Warn(LOG_TAG + "try to UnSubscribe unexist topic: " + topic)
	}
}

func (mqtt_cli *MqttClient_t) PublishMsg(topic string, data []byte) {
	logger.Fatal("unimplemented function")
	panic(nil)
}
