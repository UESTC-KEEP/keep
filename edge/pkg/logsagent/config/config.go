package config

import (
<<<<<<< HEAD
	edgeagent "keep/pkg/apis/compoenentconfig/edgeagent/v1alpha1"
=======
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
>>>>>>> b0af266029c89d24fd39eac5960a66536ae9a802
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
