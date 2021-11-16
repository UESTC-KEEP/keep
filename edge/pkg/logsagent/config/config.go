package config

import (

	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"

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
