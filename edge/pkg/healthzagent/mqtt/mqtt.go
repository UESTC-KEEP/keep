package mqtt

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"github.com/wonderivan/logger"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"keep/constants"
	_ "net/http/pprof"
	"time"
)

//config

const MqttForever = 0 //0表示无限等待

type MqtterrcodeT int

const (
	MQTT_OK MqtterrcodeT = iota
	MqttChanClosed
	MqttTimeout
	MqttTopicUnexist
)

type MqtterrretT struct {
	errCode MqtterrcodeT
}

func (err *MqtterrretT) Error() string {
	return constants.DefaultMqttLogTag //TODO 以后写点详细信息
}

type MqtttopicinfoT struct {
	dataChan    chan []byte
	timeLimitMs uint
}

type TopicMap map[string]*MqtttopicinfoT
type MqttclientT struct {
	pMqttClient *client.Client
	topicMap    TopicMap
}

type TopicConf struct {
	TopicName string
	TimeoutMs uint
}

func CreateMqttClient(clientName string, brokerIp string, brokerPort string) *MqttclientT {
	var mqttCli = client.New(&client.Options{
		ErrorHandler: func(err error) {
			logger.Error("连接mqttbroker失败...", err)
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
		logger.Fatal(err)
		panic(err)
	}

	pCli := new(MqttclientT)

	pCli.pMqttClient = mqttCli
	pCli.topicMap = make(TopicMap)

	return pCli
}

func CreateMqttClientNoName(brokerIp string, brokerPort string) *MqttclientT { //随机生成客户端名字
	return CreateMqttClient((uuid.NewV4()).String(), brokerIp, brokerPort)
}

var countr = 0

// RegistSubscribeTopic 只要调用一次就行，其后等着自己的回调函数就行，不用反复注册订阅
func (mqttCli *MqttclientT) RegistSubscribeTopic(pConf *TopicConf) {
	//hubCli:= healthzhub.NewHealzHub()
	topicName := pConf.TopicName
	pCli := mqttCli.pMqttClient
	_, exist := mqttCli.topicMap[topicName]
	if exist { //不能重复订阅同一主题
		logger.Warn(constants.DefaultMqttLogTag + ": Skip subscribeing duplicated topic " + topicName)
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
					fmt.Println(constants.DefaultMqttLogTag+": Topic= "+string(topicName)+"\tData=", message)
					// 存入sqlite进行有限制的持久化
					//hubCli.InsertIntoSqlite(message)
					mqttCli.topicMap[string(topicName)].dataChan <- message //TODO 这个地方得先设法获取通道是否满了，不然会在满了后一直阻塞
				},
			},
		},
	})
	if nil != err {
		logger.Error(constants.DefaultMqttLogTag, ":Failed to subscribe topic: ", topicName, " error=", err)
		return
	}

	mqttCli.topicMap[topicName] = &MqtttopicinfoT{
		dataChan:    make(chan []byte, constants.DefaultMqttChanSize),
		timeLimitMs: pConf.TimeoutMs,
	}
}

func (mqttCli *MqttclientT) GetTopicData(topic string) ([]byte, error) {
	topicInfo, exist := mqttCli.topicMap[topic]
	if exist {
		if MqttForever == topicInfo.timeLimitMs { //TODO 不知道怎么复用select，凑合一下
			data := <-topicInfo.dataChan
			if data == nil {
				logger.Warn(constants.DefaultMqttLogTag + ":The data channel of topic \"" + topic + "\" was closed")
				return nil, &MqtterrretT{MqttChanClosed}
			} else {
				return data, nil
			}
		} else {
			select {
			case data := <-topicInfo.dataChan:
				if data == nil {
					logger.Warn(constants.DefaultMqttLogTag + ":The data channel of topic \"" + topic + "\" was closed")
					return nil, &MqtterrretT{MqttChanClosed}
				} else {
					return data, nil
				}
			case <-time.After(time.Duration(topicInfo.timeLimitMs) * time.Millisecond):
				logger.Error(constants.DefaultMqttLogTag, ": TIMEOUT while reading topic: ", topic)
				return nil, &MqtterrretT{MqttTimeout}
			}
		}

	} else {
		logger.Error(constants.DefaultMqttLogTag + "try to read unregisted topic " + topic)
		return nil, &MqtterrretT{MqttTopicUnexist}
	}
}

func (mqttCli *MqttclientT) UnSubscribeTopic(topic string) {
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
		logger.Warn(constants.DefaultMqttLogTag + "try to UnSubscribe unexist topic: " + topic)
	}
}

func (mqttCli *MqttclientT) PublishMsg(topic string, data []byte) {
	logger.Fatal("unimplemented function")
	panic(nil)
}
