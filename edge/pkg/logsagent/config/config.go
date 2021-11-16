package config

import (
	edgeagent "keep/pkg/apis/compoenentconfig/edgeagent/v1alpha1"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	edgeagent.LogsAgent
}

func InitConfigure(l *edgeagent.LogsAgent) {
	once.Do(func() {
		Config = Configure{
			LogsAgent: *l,
		}
	})
}

func Get() *Configure {
	return &Config
}