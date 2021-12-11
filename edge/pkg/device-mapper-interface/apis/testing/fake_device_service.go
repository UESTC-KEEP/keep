package testing

import v1 "keep/edge/pkg/device-mapper-interface/apis/devices/v1"

var (
	FakeStatus = map[string]string{"temperature": "27"}
)

type FakeDeviceService struct {
}

func (fds *FakeDeviceService) GetDeviceStatus(deviceId string) (v1.DeviceStatusResponse, error) {
	return v1.DeviceStatusResponse{
		Status: &v1.DeviceStatus{
			Status: FakeStatus,
			Info:   map[string]string{"result": "ok"},
		},
		Info: map[string]string{deviceId: "ok"},
	}, nil
}
