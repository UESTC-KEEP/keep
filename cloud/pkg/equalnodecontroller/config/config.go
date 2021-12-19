package config

import (
	crdClientset "keep/cloud/pkg/client/clientset/versioned"
	"keep/cloud/pkg/common/client"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
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
		EqndClient = client.GetCRDClient()
		Config = Configure{
			EqualNodeController: *eqndc,
		}
	})
}
