package device_manage_interface

import (
	"keep/edge/pkg/common/modules"
	dmi_cfg "keep/edge/pkg/device_manage_interface/config"
	"keep/edge/pkg/edgepublisher/publisher"
	"keep/edge/pkg/healthzagent/mqtt"
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"keep/pkg/util/core"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	"keep/pkg/util/kplogger"
	logger "keep/pkg/util/loggerv1.0.1"
	"os"
	"time"
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

	// 测试查询云端设备列表
	go func() {
		time.Sleep(5 * time.Second)
		logger.Warn("开始测试查询divece 列表...")
		// 张连军：测试抄送一份到k8sclient 可注释之
		testmap := make(map[string]interface{})
		testmap["namespace"] = "default"
		msg_zlj := model.Message{
			Header: model.MessageHeader{},
			Router: model.MessageRoute{
				// 指明调用函数后  功能模块返回结果的接收模块(查询pod列表后由RequestDispatcher 下发节点)
				Source: modules.HealthzAgentModule,
				// 如果需要群发模块则填写此之段
				Group: "",
				// 以下两个内容由调用的资源模块进行解析 先定位到操作资源 在定位增删查改
				// 对资源进行的操作
				Operation: "list",
				// 资源所在路由
				Resource: "$uestc/keep/k8sclient/kubeedgeengin/devices/",
			},
			// 内容及参数由RequestDispatcher与被调用模块协商
			Content: testmap,
		}
		beehiveContext.Send(modules.EdgePublisherModule, msg_zlj)
	}()
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
