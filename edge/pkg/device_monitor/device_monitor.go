package devicemonitor

import (
	"fmt"
	"keep/edge/pkg/healthzagent/mqtt"
	"net/http"
)

const DEVICE_REG_PORT = "8085"
const HTTP_SERVER_ADDR = "localhost"

const MQTT_BROKER_PORT = "1883"
const MQTT_BROKER_ADDR = "192.168.1.40"

type DeviceMonitor struct {
	mqtt_cli *mqtt.MqttClient
	// http_server
	device_list []string //记录已经注册的设备
	// recorded_device_list [] string //TODO 已经在k8s中注册的设备，可以直接查到
}

func NewDeviceMonitor() *DeviceMonitor {
	monitor := new(DeviceMonitor)
	monitor.mqtt_cli = mqtt.CreateMqttClientNoName(MQTT_BROKER_ADDR, MQTT_BROKER_PORT)
	//TODO 只是开了客户端，没监听

	return monitor
}

func (monitor *DeviceMonitor) Destroy() {
	if monitor.mqtt_cli != nil {
		monitor.mqtt_cli.DestroyMqttClient()
	}
}

func (monitor *DeviceMonitor) Run() {
	go monitor.checkDevice()
	monitor.LocalDeviceRegistryServer()
}

func (monitor *DeviceMonitor) ServeHTTP(resp http.ResponseWriter, req *http.Request) { //  监听本机上的新mapper的注册请求
	fmt.Fprintln(resp, "TODO")
}

func (monitor *DeviceMonitor) LocalDeviceRegistryServer() {
	http.Handle("/", monitor)
	http.ListenAndServe(HTTP_SERVER_ADDR+":"+DEVICE_REG_PORT, nil)
}

func (monitor *DeviceMonitor) checkDevice() {
	if len(monitor.device_list) > 0 {
		for _, device := range monitor.device_list {
			fmt.Println(device)
			monitor.mqtt_cli.RegistSubscribeTopic(&mqtt.TopicConf{
				TopicName: TopicDeviceDataUpdate(device),
				TimeoutMs: 0,
				DataMode:  mqtt.MQTT_BLOCK_MODE,
			})
		}
	}
}
