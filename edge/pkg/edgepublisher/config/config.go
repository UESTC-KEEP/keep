package config

import (
	edgeagent "keep/pkg/apis/compoenentconfig/edgeagent/v1alpha1"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	edgeagent.EdgePublisher
}

func InitConfigure(ep *edgeagent.EdgePublisher) {
	once.Do(func() {
		Config = Configure{
			EdgePublisher: *ep,
		}
	})
}

func Get() *Configure {
	return &Config
}
