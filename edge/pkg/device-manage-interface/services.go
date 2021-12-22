package device_manage_interface

import v1 "keep/edge/pkg/device-manage-interface/apis/devices/v1"

type DeviceStatus interface {
	// GetDeviceStatus 获取设备的状态
	GetDeviceStatus(deviceId string) (v1.DeviceStatusResponse, error)
}

type DeviceService interface {
	DeviceStatus
}
