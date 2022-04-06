package config

import (
	edgeagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"

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
