package config

import (
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	edgeagent.DeviceMapperInterface
}

func InitConfigure(dmi *edgeagent.DeviceMapperInterface) {
	once.Do(func() {
		Config = Configure{
			DeviceMapperInterface: *dmi,
		}
	})
}

func Get() *Configure {
	return &Config
}
