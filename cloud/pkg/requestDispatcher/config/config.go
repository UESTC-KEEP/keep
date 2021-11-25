package config

import (
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	cloudagent.RequestDispatcher
	//KubeAPIConfig *v1alpha1.KubeAPIConfig
	Ca    []byte
	CaKey []byte
	Cert  []byte
	Key   []byte
}

func InitConfigure(r *cloudagent.RequestDispatcher) {
	once.Do(func() {
		Config = Configure{
			RequestDispatcher: *r,
		}
	})
}
