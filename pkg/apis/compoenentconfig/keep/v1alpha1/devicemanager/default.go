package v1alpha1

func NewDefaultDeviceManagerConfig() *DeviceManagerConfig {
	return &DeviceManagerConfig{
		Modules: &Modules{
			Device: &Device{
				Enable:                true,
				KindName:              "",
				DeviceControllMethods: nil,
			},
		},
	}
}
