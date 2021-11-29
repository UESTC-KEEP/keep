package config

import (
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	cloudagent.EqualNodeController
}

func InitConfigure(eqndc *cloudagent.EqualNodeController) {
	once.Do(func() {
		Config = Configure{
			EqualNodeController: *eqndc,
		}
	})
}
