package devicemonitor

import (
	"context"
	"encoding/json"
	"fmt"
	"keep/edge/pkg/device_monitor/mqtt"
	"keep/pkg/util/kplogger"
	"net/http"
	"time"
)

//这个实际由mapper调用，用于简化mapper发送消息的操作
//首次启动时会向device_monitor的指定端口报到和注册，然后用mqtt发消息
//确保发出的消息能被转为map

type MessageInterface struct {
	ctx         context.Context
	cancel      context.CancelFunc
	device_name string
	mqtt_cli    *mqtt.MqttClient
}

type MapperMessageDef map[string]interface{}
type MapperMessage struct {
	data MapperMessageDef
}

func NewMsgInterface(device_name string) *MessageInterface {
	msg_interface := new(MessageInterface)

	ctx, cancel := context.WithCancel(context.Background())
	msg_interface.ctx = ctx
	msg_interface.cancel = cancel

	msg_interface.mqtt_cli = mqtt.CreateMqttClientNoName(MQTT_BROKER_ADDR, MQTT_BROKER_PORT)
	msg_interface.device_name = device_name

	go msg_interface.tryRegistToDeviceMonitor()

	go msg_interface.deviceNameReporterDaemon()

	return msg_interface
}

func (obj *MessageInterface) tryRegistToDeviceMonitor() {
	retry_period := 10 * time.Second
	retry_timer := time.NewTimer(retry_period) //time.After会溢出
	defer retry_timer.Stop()
	for {
		err := obj.registToDeviceMonitor()
		if err == nil { //每10s尝试一次，成功后就停止
			return
		}
		retry_timer.Reset(retry_period)
		select {
		case <-obj.ctx.Done():
		case <-retry_timer.C: //TODO
		}
	}
}

func (obj *MessageInterface) registToDeviceMonitor() error { //目前只是把本设备名称通知给device monitor
	kplogger.Infof("Try to Regist device <%s> to DM", obj.device_name)
	url := "http://" + HTTP_SERVER_ADDR + ":" + DEVICE_REG_PORT + "/" + obj.device_name
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		kplogger.Error("Cannont create http req")
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		kplogger.Errorf("Cannont regist device <%s>", obj.device_name)
		fmt.Println(err)
		return err
	}
	return nil
}

func (obj *MessageInterface) Destroy() {
	obj.cancel()
	if obj.mqtt_cli != nil {
		obj.mqtt_cli.DestroyMqttClient()
	}
}

//额外添加处理DM广播设备发现的接口 ,收到DM发的广播后，就会向DM报告本设备的名称
func (obj *MessageInterface) deviceNameReporterDaemon() {
	topic := TopicInquireDeviceName()
	obj.mqtt_cli.RegistSubscribeTopic(&mqtt.TopicConf{
		TopicName: topic,
		TimeoutMs: 0,
		DataMode:  mqtt.MQTT_BLOCK_MODE,
	})

	for { //mqtt客户端在关闭后会结束阻塞
		_, err := obj.mqtt_cli.GetTopicData(topic) //TODO 应该写点数据，当作校验
		if nil != err {
			kplogger.Error(err)
			continue
		}
		obj.registToDeviceMonitor()
	}
}

func NewMapperMsg() *MapperMessage {
	msg := new(MapperMessage)
	msg.data = make(MapperMessageDef)
	//BaseMessage
	msg.data["EventID"] = "0" //TODO 这两个地方是随便填的
	msg.data["Timestamp"] = 123456789
	return msg
}

func (msg *MapperMessage) AddItem(item_name string, data interface{}) {
	msg.data[item_name] = data
}

//TODO 还要实现其他的mapper和edgetopic接口
func (obj *MessageInterface) SendStatusData(msg *MapperMessage) {
	topic := TopicDeviceDataUpdate(obj.device_name)
	mjson, _ := json.Marshal(msg.data) //用map生成json
	obj.mqtt_cli.PublishMsg(topic, mjson, false)
}
