package config

import (
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	cloudagent.CloudImageManager
}

func InitConfigure(cmi *cloudagent.CloudImageManager) {
	once.Do(func() {
		Config = Configure{
			CloudImageManager: *cmi,
		}
	})
}

func Get() *Configure {
	return &Config
}
