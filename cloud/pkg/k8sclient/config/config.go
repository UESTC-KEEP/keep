package config

import (
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	cloudagent.K8sClient
}

func InitConfigure(k *cloudagent.K8sClient) {
	once.Do(func() {
		Config = Configure{
			K8sClient: *k,
		}
	})
}

func Get() *Configure {
	return &Config
}
