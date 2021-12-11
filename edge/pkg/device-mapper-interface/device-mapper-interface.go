package device_mapper_interface

import (
	"keep/edge/pkg/common/modules"
	devicemapperinterfaceconfig "keep/edge/pkg/device-mapper-interface/config"
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"keep/pkg/util/core"
	logger "keep/pkg/util/loggerv1.0.1"
	"os"
)

// 我叫interface但是我不是一个interface 我是一个模块

type DeviceMapperInterface struct {
	enable bool `json:"enable"`
}

// Register 注册模块
func Register(dmi *edgeagent.DeviceMapperInterface) {
	devicemapperinterfaceconfig.InitConfigure(dmi)
	devicemapperinter, err := NewDeviceMapperInterface(dmi.Enable)
	if err != nil {
		logger.Error("初始化DeviceMapperInterface失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(devicemapperinter)
}

func (dmi *DeviceMapperInterface) Cleanup() {}

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

}

func NewDeviceMapperInterface(enable bool) (*DeviceMapperInterface, error) {
	return &DeviceMapperInterface{
		enable: enable,
	}, nil
}
