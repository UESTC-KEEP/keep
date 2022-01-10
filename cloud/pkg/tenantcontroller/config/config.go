package config

import (
	cloudagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	cloudagent.TenantController
}

func InitConfigure(tc *cloudagent.TenantController) {
	once.Do(func() {
		Config = Configure{
			TenantController: *tc,
		}
	})
}

func Get() *Configure {
	return &Config
}
