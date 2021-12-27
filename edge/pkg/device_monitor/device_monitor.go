package devicemonitor

import (
	"context"
	"fmt"
	"keep/edge/pkg/healthzagent/mqtt"
	"keep/pkg/util/kplogger"
	"net/http"
	"time"
)

const DEVICE_REG_PORT = "8085"
const HTTP_SERVER_ADDR = "localhost"

const MQTT_BROKER_PORT = "1883"
const MQTT_BROKER_ADDR = "192.168.1.40"

type DeviceMonitor struct {
	ctx      context.Context
	cancel   context.CancelFunc
	mqtt_cli *mqtt.MqttClient
	// http_server
	device_list map[string]interface{} //记录已经注册的设备
	// recorded_device_list [] string //TODO 已经在k8s中注册的设备，可以直接查到
}

func NewDeviceMonitor() *DeviceMonitor {
	monitor := new(DeviceMonitor)

	ctx, cancel := context.WithCancel(context.Background())
	monitor.ctx = ctx
	monitor.cancel = cancel
	monitor.device_list = make(map[string]interface{})
	monitor.mqtt_cli = mqtt.CreateMqttClientNoName(MQTT_BROKER_ADDR, MQTT_BROKER_PORT)
	//TODO 只是开了客户端，没监听

	return monitor
}

func (monitor *DeviceMonitor) Destroy() {
	monitor.cancel()
	if monitor.mqtt_cli != nil {
		monitor.mqtt_cli.DestroyMqttClient()
	}
}

func (monitor *DeviceMonitor) Run() {
	go monitor.monitorDevice()
	monitor.LocalDeviceRegistryServer()
}

func (monitor *DeviceMonitor) ServeHTTP(resp http.ResponseWriter, req *http.Request) { //  监听本机上的新mapper的注册请求
	fmt.Fprintln(resp, req.URL.String(), "TODO")
	device_name := req.URL.String()
	monitor.addDeviceToRecord(device_name[1:]) //去掉第一个斜杠
}

func (monitor *DeviceMonitor) LocalDeviceRegistryServer() {
	http.Handle("/", monitor)
	http.ListenAndServe(HTTP_SERVER_ADDR+":"+DEVICE_REG_PORT, nil)
}

func (monitor *DeviceMonitor) monitorDevice() {
	if len(monitor.device_list) > 0 {
		for device, _ := range monitor.device_list {
			// fmt.Println(device)
			monitor.mqtt_cli.RegistSubscribeTopic(&mqtt.TopicConf{
				TopicName: TopicDeviceDataUpdate(device),
				TimeoutMs: 0,
				DataMode:  mqtt.MQTT_BLOCK_MODE,
			})
			time.Sleep(time.Millisecond * 500) //TODO 暂定
		}
	}
}

func (monitor *DeviceMonitor) InquireLocalDevice() {
	topic := TopicInquireDeviceName()
	monitor.mqtt_cli.PublishMsg(topic, []byte("INQUIRE"), false) // TODO 内容暂定，好像不要也行
}

func (monitor *DeviceMonitor) addDeviceToRecord(device_name string) {
	_, exist := monitor.device_list[device_name]
	if exist {
		kplogger.Warnf("Device <%s> is already in the record", device_name)
	} else {
		monitor.device_list[device_name] = "TODO" //TODO 以后加更多的信息
		kplogger.Infof("Add device <%s> to record", device_name)

	}
}

func (monitor *DeviceMonitor) RemoveDevice(device_name string) {
	_, exist := monitor.device_list[device_name]
	if exist {
		delete(monitor.device_list, device_name)
		kplogger.Info("Remove device <%s> from record", device_name)
	} else {
		kplogger.Warnf("Device <%s> is not in the record", device_name)
	}
}
