package config

import (
	cloudagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	cloudagent.PromServer
}

func InitConfigure(ps *cloudagent.PromServer) {
	once.Do(func() {
		Config = Configure{
			PromServer: *ps,
		}
	})
}

func Get() *Configure {
	return &Config
}
