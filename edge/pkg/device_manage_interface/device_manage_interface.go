package device_manage_interface

import (
	"keep/edge/pkg/common/modules"
	dmi_cfg "keep/edge/pkg/device_manage_interface/config"
	"keep/edge/pkg/device_monitor/mqtt"
	"keep/edge/pkg/edgepublisher/publisher"

	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"keep/pkg/util/core"
	"keep/pkg/util/core/model"
	"keep/pkg/util/kplogger"
	"os"
)

// 我叫interface但是我不是一个interface 我是一个模块

type DeviceManageInterface struct {
	mqtt_cli *mqtt.MqttClient
	enable   bool `json:"enable"`
}

// Register 注册模块
func Register(dmi *edgeagent.DeviceMapperInterface) {
	dmi_cfg.InitConfigure(dmi)
	device_manage_interface, err := NewDeviceManageInterface(dmi.Enable)
	if err != nil {
		kplogger.Error("初始化DeviceMapperInterface失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(device_manage_interface)
}

func (dmi *DeviceManageInterface) Cleanup() {
	if !(dmi.enable) {
		return
	}

	if nil != dmi.mqtt_cli { // TODO
		dmi.mqtt_cli.DestroyMqttClient()
	}

}

func (dmi *DeviceManageInterface) Group() string {
	return modules.DeviceMapperInterfaceGroup
}

func (dmi *DeviceManageInterface) Name() string {
	return modules.DeviceMapperInterfaceModule
}

func (dmi *DeviceManageInterface) Enable() bool {
	return dmi.enable
}

func (dmi *DeviceManageInterface) Start() {
	// 获取设备列表
	var msg model.Message
	msg.SetResourceOperation("$uestc/keep/k8sclient/kubeedgeengin/devices/", "list")

	publisher.Publish(msg)
}

func NewDeviceManageInterface(enable bool) (*DeviceManageInterface, error) {
	dmi_obj := new(DeviceManageInterface)
	dmi_obj.mqtt_cli = mqtt.CreateMqttClientNoName("localhost", "1883")
	if nil == dmi_obj.mqtt_cli {
		dmi_obj.enable = false
	} else {
		dmi_obj.enable = enable
	}

	return dmi_obj, nil
}
