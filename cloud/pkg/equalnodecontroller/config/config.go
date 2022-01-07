package config

import (
	crdClientset "github.com/UESTC-KEEP/keep/cloud/pkg/client/eqnd/clientset/versioned"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	cloudagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once
var EqndClient crdClientset.Interface

type Configure struct {
	cloudagent.EqualNodeController
}

func InitConfigure(eqndc *cloudagent.EqualNodeController) {
	once.Do(func() {
		EqndClient = client.GetEqndCRDClient()
		Config = Configure{
			EqualNodeController: *eqndc,
		}
	})
}
