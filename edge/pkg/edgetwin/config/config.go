package config

import (
	edgeagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"

	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	edgeagent.EdgeTwin
}

func InitConfigure(et *edgeagent.EdgeTwin) {
	once.Do(func() {
		Config = Configure{
			EdgeTwin: *et,
		}
	})
}

func Get() *Configure {
	return &Config
}
