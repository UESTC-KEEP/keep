package device_manage_interface

import (
	"keep/edge/pkg/common/modules"
	devicemapperinterfaceconfig "keep/edge/pkg/device_manage_interface/config"
	"keep/edge/pkg/healthzagent/mqtt"
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"keep/pkg/util/core"
	"keep/pkg/util/kplogger"
	"os"
)

// 我叫interface但是我不是一个interface 我是一个模块

type DeviceMapperInterface struct {
	mqtt_cli *mqtt.MqttClient
	enable   bool `json:"enable"`
}

// Register 注册模块
func Register(dmi *edgeagent.DeviceMapperInterface) {
	devicemapperinterfaceconfig.InitConfigure(dmi)
	devicemapperinter, err := NewDeviceMapperInterface(dmi.Enable)
	if err != nil {
		kplogger.Error("初始化DeviceMapperInterface失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(devicemapperinter)
}

func (dmi *DeviceMapperInterface) Cleanup() {
	if !(dmi.enable) {
		return
	}

	if nil != dmi.mqtt_cli { // TODO
		dmi.mqtt_cli.DestroyMqttClient()
	}

}

func (dmi *DeviceMapperInterface) Group() string {
	return modules.DeviceMapperInterfaceGroup
}

func (dmi *DeviceMapperInterface) Name() string {
	return modules.DeviceMapperInterfaceModule
}

func (dmi *DeviceMapperInterface) Enable() bool {
	return dmi.enable
}

func (dmi *DeviceMapperInterface) Start() {
	//获取设备列表
	// var msg model.Message
	// msg.Router
	// publisher.Publish(msg)
}

func NewDeviceMapperInterface(enable bool) (*DeviceMapperInterface, error) {
	dmi_obj := new(DeviceMapperInterface)
	dmi_obj.mqtt_cli = mqtt.CreateMqttClientNoName("192.168.1.40", "1883")
	if nil == dmi_obj.mqtt_cli {
		dmi_obj.enable = false
	} else {
		dmi_obj.enable = enable
	}

	return dmi_obj, nil
}
