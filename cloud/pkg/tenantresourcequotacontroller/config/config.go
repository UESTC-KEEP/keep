package config

import (
	trqcrdClientset "github.com/UESTC-KEEP/keep/cloud/pkg/client/trq/clientset/versioned"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	cloudagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
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
