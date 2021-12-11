package options

import (
	"keep/constants/device"
	v1alpha1 "keep/pkg/apis/compoenentconfig/keep/v1alpha1/devicemanager"
	logger "keep/pkg/util/loggerv1.0.1"
)

/*
Copyright 2021 The KEEP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

type DeviceManagerOptions struct {
	ConfigFile string
}

func NewDefaultDeviceManagerConfig() *DeviceManagerOptions {
	return &DeviceManagerOptions{
		ConfigFile: device.DefaultDeviceConfigFile,
	}
}

func (o *DeviceManagerOptions) Config() (*v1alpha1.DeviceManagerConfig, error) {
	cfg := v1alpha1.NewDefaultDeviceManagerConfig()
	// 先读默认配置  然后在从配置文件读取用户配置   以用户配置的配置文件为主
	if err := cfg.Parse(o.ConfigFile); err != nil {
		logger.Error(err)
		return nil, err
	}
	return cfg, nil
}
