package config

import (
	crdClientset "keep/cloud/pkg/client/clientset/versioned"
	"keep/cloud/pkg/common/client"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once
var KeepCrdClient crdClientset.Interface

type Configure struct {
	cloudagent.TenantResourceQuotaController
}

func InitConfigure(trq *cloudagent.TenantResourceQuotaController) {
	once.Do(func() {
		KeepCrdClient = client.GetCRDClient()
		Config = Configure{
			TenantResourceQuotaController: *trq,
		}
	})
}
