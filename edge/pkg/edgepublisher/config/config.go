package config

import (
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
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
