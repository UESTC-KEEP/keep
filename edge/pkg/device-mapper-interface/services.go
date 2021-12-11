package device_mapper_interface

import "keep/edge/pkg/device-mapper-interface/apis/devices/v1"

type DeviceStatus interface {
	// GetDeviceStatus 获取设备的状态
	GetDeviceStatus(deviceId string) (v1.DeviceStatusResponse, error)
}

type DeviceService interface {
	DeviceStatus
}
