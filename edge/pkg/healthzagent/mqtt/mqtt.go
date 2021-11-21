package mqtt

import (
	"fmt"
	"keep/pkg/util/kelogger"
	"time"

	_ "github.com/mattn/go-sqlite3"

	uuid "github.com/satori/go.uuid"

	"keep/constants"
	_ "net/http/pprof"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

//config

const MqttForever = 0 //0表示无限等待

type MqttErrCode int

const (
	MQTT_OK MqttErrCode = iota
	MQTT_CHAN_CLOSED
	MQTT_TIME_OUT
	MQTT_TOPIC_UNEXIST
)

type MqttErrRet struct {
	errCode MqttErrCode
}

func (err *MqttErrRet) Error() string {
	return constants.DefaultMqttLogTag //TODO 以后写点详细信息
}

type mqttTopicInfo struct { //不导出这个结构体
	dataChan    chan []byte
	timeLimitMs uint
}

type TopicMap map[string]*mqttTopicInfo
type MqttClient struct {
	pMqttClient *client.Client
	topicMap    TopicMap
}

type TopicConf struct {
	TopicName string
	TimeoutMs uint
}

func CreateMqttClient(clientName string, brokerIp string, brokerPort string) *MqttClient {
	var mqttCli = client.New(&client.Options{
		ErrorHandler: func(err error) {
			kelogger.Error("连接mqttbroker失败...", err)
		},
	})
	defer mqttCli.Terminate()

	connOpt := client.ConnectOptions{
		Network:  "tcp",
		Address:  brokerIp + ":" + brokerPort,
		ClientID: []byte(clientName),
	}
	err := mqttCli.Connect(&connOpt)

	if nil != err {
		kelogger.Fatal(err)
		panic(err)
	}

	pCli := new(MqttClient)

	pCli.pMqttClient = mqttCli
	pCli.topicMap = make(TopicMap)

	return pCli
}

func CreateMqttClientNoName(brokerIp string, brokerPort string) *MqttClient { //随机生成客户端名字
	return CreateMqttClient((uuid.NewV4()).String(), brokerIp, brokerPort)
}

var countr = 0

// RegistSubscribeTopic 只要调用一次就行，其后等着自己的回调函数就行，不用反复注册订阅
func (mqttCli *MqttClient) RegistSubscribeTopic(pConf *TopicConf) {
	//hubCli:= healthzhub.NewHealzHub()
	topicName := pConf.TopicName
	pCli := mqttCli.pMqttClient
	_, exist := mqttCli.topicMap[topicName]
	if exist { //不能重复订阅同一主题
		kelogger.Warn(constants.DefaultMqttLogTag + ": Skip subscribeing duplicated topic " + topicName)
		return
	}

	err := pCli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			{
				TopicFilter: []byte(topicName),
				QoS:         mqtt.QoS0,
				Handler: func(topicName, message []byte) {

					fmt.Println(countr)
					countr++
					kelogger.Trace(constants.DefaultMqttLogTag+": Topic= "+string(topicName)+"\tData=", message)
					// 存入sqlite进行有限制的持久化
					//hubCli.InsertIntoSqlite(message)
					mqttCli.topicMap[string(topicName)].dataChan <- message //TODO 这个地方得先设法获取通道是否满了，不然会在满了后一直阻塞

				},
			},
		},
	})
	if nil != err {
		kelogger.Error(constants.DefaultMqttLogTag, ":Failed to subscribe topic: ", topicName, " error=", err)
		return
	}

	mqttCli.topicMap[topicName] = &mqttTopicInfo{
		dataChan:    make(chan []byte, constants.DefaultMqttChanSize),
		timeLimitMs: pConf.TimeoutMs,
	}
}

func (mqttCli *MqttClient) GetTopicData(topic string) ([]byte, error) {
	topicInfo, exist := mqttCli.topicMap[topic]
	if exist {
		if MqttForever == topicInfo.timeLimitMs { //TODO 不知道怎么复用select，凑合一下
			data := <-topicInfo.dataChan
			if data == nil {
				kelogger.Warn(constants.DefaultMqttLogTag + ":The data channel of topic \"" + topic + "\" was closed")
				return nil, &MqttErrRet{MQTT_CHAN_CLOSED}
			} else {
				return data, nil
			}
		} else {
			select {
			case data := <-topicInfo.dataChan:
				if data == nil {
					kelogger.Warn(constants.DefaultMqttLogTag + ":The data channel of topic \"" + topic + "\" was closed")
					return nil, &MqttErrRet{MQTT_CHAN_CLOSED}
				} else {
					return data, nil
				}
			case <-time.After(time.Duration(topicInfo.timeLimitMs) * time.Millisecond):
				kelogger.Error(constants.DefaultMqttLogTag, ": TIMEOUT while reading topic: ", topic)
				return nil, &MqttErrRet{MQTT_TIME_OUT}
			}
		}

	} else {
		kelogger.Error(constants.DefaultMqttLogTag + "try to read unregisted topic " + topic)
		return nil, &MqttErrRet{MQTT_TOPIC_UNEXIST}
	}
}

func (mqttCli *MqttClient) UnSubscribeTopic(topic string) {
	topicInfo, exist := mqttCli.topicMap[topic]
	if exist {
		mqttCli.pMqttClient.Unsubscribe(
			&client.UnsubscribeOptions{
				TopicFilters: [][]byte{
					[]byte(topic)},
			})
		close(topicInfo.dataChan)
		//chan 关闭时的原则是：不要在接收协程中关闭，并且，如果有多个发送者时就不要关闭chan了。
		//https://studygolang.com/articles/9478
		delete(mqttCli.topicMap, topic)
	} else {
		kelogger.Warn(constants.DefaultMqttLogTag + "try to UnSubscribe unexist topic: " + topic)
	}
}

func (mqtt_cli *MqttClient) PublishMsg(topic string, data []byte) {
	kelogger.Fatal("unimplemented function")
	panic(nil)
}
