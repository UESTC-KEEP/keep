package v1alpha1

type Modules struct {
	Device *Device `json:"device,omitempty"`
}

type DeviceManagerConfig struct {
	// edgeagentagent中的模块配置
	Modules *Modules `json:"modules,omitempty"`
}

type Device struct {
	Enable   bool   `json:"enable"`
	KindName string `json:"kindName,omitempty"`
	// 写调用方式的方式路径
	DeviceControllMethods map[string]string `json:"deviceControllMethods,omitempty"`
}
