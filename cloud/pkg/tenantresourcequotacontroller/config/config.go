package config

import (
	trqcrdClientset "keep/cloud/pkg/client/trq/clientset/versioned"
	"keep/cloud/pkg/common/client"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once
var KeepCrdClient trqcrdClientset.Interface

type Configure struct {
	cloudagent.TenantResourceQuotaController
}

func InitConfigure(trq *cloudagent.TenantResourceQuotaController) {
	once.Do(func() {
		KeepCrdClient = client.GetTrqCRDClient()
		Config = Configure{
			TenantResourceQuotaController: *trq,
		}
	})
}
