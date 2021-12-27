package devicemonitor

import (
	"keep/edge/pkg/healthzagent/mqtt"
	"keep/pkg/util/kplogger"
	"net/http"
)

//这个实际由mapper调用，用于简化mapper发送消息的操作
//首次启动时会向device_monitor的指定端口报到和注册，然后用mqtt发消息
//确保发出的消息能被转为map

type MessageInterface struct {
	device_name string
	mqtt_cli    *mqtt.MqttClient
}

func NewMsgInterface(device_name string) *MessageInterface {
	msg_interface := new(MessageInterface)
	msg_interface.mqtt_cli = mqtt.CreateMqttClientNoName(MQTT_BROKER_ADDR, MQTT_BROKER_PORT)
	msg_interface.device_name = device_name

	msg_interface.registToDeviceMonitor()

	go deviceNameReporter()

	return msg_interface
}

func (obj *MessageInterface) registToDeviceMonitor() { //目前只是把本设备名称通知给device monitor
	url := HTTP_SERVER_ADDR + ":" + DEVICE_REG_PORT + "/" + obj.device_name
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		kplogger.Emer("Cannont regist device =", obj.device_name)
		return
	}

	client := &http.Client{}
	client.Do(req)
}

func (obj *MessageInterface) Destroy() {
	if obj.mqtt_cli != nil {
		obj.mqtt_cli.DestroyMqttClient()
	}
}

//TODO 还要实现其他的mapper和edgetopic接口

func (obj *MessageInterface) SendStatusData(data []byte) {
	topic := TopicDeviceDataUpdate(obj.device_name)
	obj.mqtt_cli.PublishMsg(topic, data)
}

//额外添加处理DM广播设备发现的接口 ,收到DM发的广播后，就会向DM报告本设备的名称
func (obj *MessageInterface) deviceNameReporter() {
	topic := TopicInquireDeviceName()
	obj.mqtt_cli.RegistSubscribeTopic(&mqtt.TopicConf{
		TopicName: topic,
		TimeoutMs: 0,
		DataMode:  mqtt.MQTT_BLOCK_MODE,
	})

	for {
		_, err := obj.mqtt_cli.GetTopicData(topic) //TODO 应该写点数据，当作校验
		if nil != err {
			kplogger.Error(err)
		}
		obj.registToDeviceMonitor()
	}
}
