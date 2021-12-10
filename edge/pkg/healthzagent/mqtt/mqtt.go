package mqtt

import (
	"keep/constants/edge"
	"keep/pkg/util/kplogger"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"

	uuid "github.com/satori/go.uuid"

	_ "net/http/pprof"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

//config

const MqttForever = 0 //0表示无限等待

type MqttErrCode uint8

const (
	MQTT_OK MqttErrCode = iota
	MQTT_CHAN_CLOSED
	MQTT_TIME_OUT
	MQTT_TOPIC_UNEXIST
	MQTT_SYS_ERROR
	MQTT_NO_DATA // 缓存模式下可能存在首次收到设备的消息前就杜数据的情况此时用is_visited标记是否写过数据
)

type MqttErrRet struct {
	errCode MqttErrCode
}

func (err *MqttErrRet) Error() string {
	switch err.errCode {
	case MQTT_CHAN_CLOSED:
		return "MQTT_CHAN_CLOSED" //TODO 以后写点详细信息
	case MQTT_TIME_OUT:
		return "MQTT_TIME_OUT"
	case MQTT_TOPIC_UNEXIST:
		return "MQTT_TOPIC_UNEXIST"
	case MQTT_SYS_ERROR:
		return "MQTT_SYS_ERROR"
	case MQTT_NO_DATA:
		return "MQTT_NO_DATA"
	default:
		return "MQTT_UNDEFIEND_ERROR"
	}

}

type MqttDataMode_t uint8

const (
	MQTT_BLOCK_MODE MqttDataMode_t = iota //阻塞，直到读取最新消息
	MQTT_CACHE_MODE                       //读取最新缓存的消息，不一定是当前时刻的消息
)

type TopicConf struct {
	TopicName string
	TimeoutMs int64
	DataMode  MqttDataMode_t //缓存模式下的超时检测比较麻烦
}

type mqttCachedData_t struct { //缓存模式下需要记录收集到数据的时间
	data_cache []byte       //TODO 这地方是存指针还是存值？
	cache_lock sync.RWMutex //TODO 以后试试用Mutex时的性能
	time_stamp int64
	is_visited bool //标记是否真有设备写过数据
}

type mqttBlockedData_t struct {
	stop_chan chan struct{}
	data_chan chan []byte
}

type mqttTopicInfo struct {
	data        interface{}
	timeLimitMs int64          //只限为毫秒级
	dataMode    MqttDataMode_t //TODO 这个有字节对齐的问题吗?
}

type TopicMap_t map[string]*mqttTopicInfo
type MqttClient struct {
	pMqttClient *client.Client
	topicMap    TopicMap_t
}

func CreateMqttClient(clientName string, broker_ip string, brokerPort string) *MqttClient {
	var mqttCli = client.New(&client.Options{
		ErrorHandler: func(err error) {
			kplogger.Error("连接mqttbroker失败...", err)
		},
	})

	connOpt := client.ConnectOptions{
		Network:  "tcp",
		Address:  broker_ip + ":" + brokerPort,
		ClientID: []byte(clientName),
	}
	err := mqttCli.Connect(&connOpt)

	if nil != err {
		kplogger.Fatal(err)
		panic(err)
	}

	pCli := new(MqttClient)

	pCli.pMqttClient = mqttCli
	pCli.topicMap = make(TopicMap_t)

	return pCli
}

func CreateMqttClientNoName(broker_ip string, brokerPort string) *MqttClient { //随机生成客户端名字
	return CreateMqttClient((uuid.NewV4()).String(), broker_ip, brokerPort)
}

func (mqttCli *MqttClient) DestroyMqttClient() {
	err := mqttCli.pMqttClient.Disconnect()
	if nil == err {
		kplogger.Infof("%s: mqtt client disconnected", edge.DefaultMqttLogTag)
	} else {
		kplogger.Errorf("%s: Error occured while disconnecting mqtt client", edge.DefaultMqttLogTag)
	}

	mqttCli.pMqttClient.Terminate()
	//通过Go的内部函数mapclear方法删除。这个函数并没有显示的调用方法，
	//当使用for循环遍历删除所有元素时，Go的编译器会优化成Go内部函数mapclear。
	for topic, topicInfo := range mqttCli.topicMap {

		if topicInfo.dataMode == MQTT_BLOCK_MODE {
			data_b, _ := topicInfo.data.(*mqttBlockedData_t)
			data_b.stop_chan <- struct{}{}
		}

		delete(mqttCli.topicMap, topic)
	}
}

func (mqttCli *MqttClient) clientReceivehandle(topicName, message []byte) {
	kplogger.Trace(edge.DefaultMqttLogTag+": Topic= "+string(topicName)+"\tData=", message)

	p_data_tmp := &(mqttCli.topicMap[string(topicName)].data)
	switch (*p_data_tmp).(type) { //类型断言
	case *mqttBlockedData_t:
		data_val, _ := (*p_data_tmp).(*mqttBlockedData_t) //FIXME  这个地方是怎么引用的？看起来是指针引用
		data_val.data_chan <- message                     //TODO 这个地方得先设法获取通道是否满了，不然会在满了后一直阻塞
	case *mqttCachedData_t:
		data_val, _ := (*p_data_tmp).(*mqttCachedData_t)
		data_val.cache_lock.Lock()
		defer data_val.cache_lock.Unlock()

		data_val.time_stamp = time.Now().UnixMilli() //写消息时附带记录时间，Get读取时会用当前时间和记录的时间做对比，判断是否超时
		data_val.data_cache = message                //FIXME 这里应该有多线程同步的问题
		data_val.is_visited = true

	}
}

// RegistSubscribeTopic 只要调用一次就行，其后等着自己的回调函数就行，不用反复注册订阅
func (mqttCli *MqttClient) RegistSubscribeTopic(pConf *TopicConf) {
	//hubCli:= healthzhub.NewHealzHub()
	topic_name := pConf.TopicName
	pCli := mqttCli.pMqttClient
	_, exist := mqttCli.topicMap[topic_name]
	if exist { //不能重复订阅同一主题
		kplogger.Warn(edge.DefaultMqttLogTag + ": Skip subscribeing duplicated topic " + topic_name)
		return
	}

	sub_opt := client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			{
				TopicFilter: []byte(topic_name),
				QoS:         mqtt.QoS0,
				Handler:     mqttCli.clientReceivehandle,
			},
		},
	}

	mqttCli.topicMap[topic_name] = &mqttTopicInfo{ //TODO map得考虑多线程互斥问题
		timeLimitMs: pConf.TimeoutMs,
		dataMode:    pConf.DataMode,
	}
	switch pConf.DataMode {
	case MQTT_BLOCK_MODE:
		mqttCli.topicMap[topic_name].data = &mqttBlockedData_t{ //以后通过类型断言做判断
			stop_chan: make(chan struct{}),
			data_chan: make(chan []byte, edge.DefaultMqttChanSize),
		}
	case MQTT_CACHE_MODE:
		mqttCli.topicMap[topic_name].data = &mqttCachedData_t{
			data_cache: nil,
			time_stamp: time.Now().UnixMilli(),
			is_visited: false,
		}
	}

	err := pCli.Subscribe(&sub_opt)
	if nil != err {
		delete(mqttCli.topicMap, topic_name)
		kplogger.Error(edge.DefaultMqttLogTag, ":Failed to subscribe topic: ", topic_name, " error=", err)
		return
	}
}

func (mqttCli *MqttClient) getDataBlockMode(topic string, topicInfo *mqttTopicInfo) ([]byte, error) {
	data_b, is_block_mode := topicInfo.data.(*mqttBlockedData_t)
	if !is_block_mode {
		return nil, &MqttErrRet{MQTT_SYS_ERROR}
	}

	if MqttForever == topicInfo.timeLimitMs { //TODO 不知道怎么复用select，凑合一下
		select {
		case data := <-data_b.data_chan:
			if data == nil {
				kplogger.Warn(edge.DefaultMqttLogTag + ":The data channel of topic \"" + topic + "\" was closed")
				return nil, &MqttErrRet{MQTT_CHAN_CLOSED}
			} else {
				return data, nil
			}
		case <-data_b.stop_chan:
			kplogger.Warn(edge.DefaultMqttLogTag + ":The data channel of topic \"" + topic + "\" was closed")
			return nil, &MqttErrRet{MQTT_CHAN_CLOSED}
		}

	} else {
		select {
		case data := <-data_b.data_chan:
			if data == nil {
				kplogger.Warn(edge.DefaultMqttLogTag + ":The data channel of topic \"" + topic + "\" was closed")
				return nil, &MqttErrRet{MQTT_CHAN_CLOSED}
			} else {
				return data, nil
			}
		case <-time.After(time.Duration(topicInfo.timeLimitMs) * time.Millisecond):
			kplogger.Error(edge.DefaultMqttLogTag, ": TIMEOUT while reading topic: ", topic)
			return nil, &MqttErrRet{MQTT_TIME_OUT}
		case <-data_b.stop_chan:
			kplogger.Warn(edge.DefaultMqttLogTag + ":The data channel of topic \"" + topic + "\" was closed")
			return nil, &MqttErrRet{MQTT_CHAN_CLOSED}
		}
	}
}

func (mqttCli *MqttClient) getDataCacheMode(topic string, topicInfo *mqttTopicInfo) ([]byte, error) {
	//读取时用系统当前时间和记录的时间做比较，判断是否超时，如果没超时，就返回最新缓存的数据
	data_c, is_cache_mode := topicInfo.data.(*mqttCachedData_t)
	if !is_cache_mode { //这里 校验不算多余，因为写错类型时，interface不会报错
		return nil, &MqttErrRet{MQTT_SYS_ERROR}
	}

	data_c.cache_lock.RLock()
	defer data_c.cache_lock.RUnlock()

	if !(data_c.is_visited) { //也许以后用nil就行，让外面读数据的去判断数据是否有效？TODO
		kplogger.Error(edge.DefaultMqttLogTag, ":  topic: ", topic, "  no data")
		return nil, &MqttErrRet{MQTT_NO_DATA} //这个地方就只有开始那一瞬间用用处了
	}

	if topicInfo.timeLimitMs < time.Now().UnixMilli()-data_c.time_stamp {
		kplogger.Error(edge.DefaultMqttLogTag, ": TIMEOUT while reading topic: ", topic)
		return nil, &MqttErrRet{MQTT_TIME_OUT}
	}
	data_ret := data_c.data_cache

	return data_ret, nil
}

func (mqttCli *MqttClient) GetTopicData(topic string) ([]byte, error) {
	topicInfo, exist := mqttCli.topicMap[topic]
	if exist {
		switch topicInfo.dataMode {
		case MQTT_BLOCK_MODE:
			return mqttCli.getDataBlockMode(topic, topicInfo)
		case MQTT_CACHE_MODE:
			return mqttCli.getDataCacheMode(topic, topicInfo)
		default:
			return nil, &MqttErrRet{MQTT_SYS_ERROR}
		}
	} else {
		kplogger.Error(edge.DefaultMqttLogTag + "try to read unregisted topic " + topic)
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
		switch topicInfo.data.(type) {
		case mqttBlockedData_t:
			data_b, is_block_mode := topicInfo.data.(*mqttBlockedData_t)
			if is_block_mode {
				close(data_b.data_chan)
			} else {
				kplogger.Errorf("%s,data mode error", edge.DefaultMqttLogTag)
			}

		case mqttCachedData_t: //FIXME 可以直接什么都不做？golang没free？
		}

		//chan 关闭时的原则是：不要在接收协程中关闭，并且，如果有多个发送者时就不要关闭chan了。
		//https://studygolang.com/articles/9478
		delete(mqttCli.topicMap, topic)
	} else {
		kplogger.Warn(edge.DefaultMqttLogTag + "try to UnSubscribe unexist topic: " + topic)
	}
}

func (mqtt_cli *MqttClient) GetTopicNum() int { //返回当前监听的topic数目
	return len(mqtt_cli.topicMap)
}

func (mqtt_cli *MqttClient) PublishMsg(topic string, data []byte) {
	kplogger.Fatal("unimplemented function")
	panic(nil)
}
