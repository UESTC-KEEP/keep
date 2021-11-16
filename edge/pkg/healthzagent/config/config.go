package config

import (
	edgeagent "keep/pkg/apis/compoenentconfig/edgeagent/v1alpha1"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	edgeagent.HealthzAgent
}

func InitConfigure(h *edgeagent.HealthzAgent) {
	once.Do(func() {
		Config = Configure{
			HealthzAgent: *h,
		}
	})
}

func Get() *Configure {
	return &Config
}
