package equalnodecontroller

import (
	flag "github.com/spf13/pflag"
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/equalnodecontroller/config"
	"keep/cloud/pkg/equalnodecontroller/controller"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/core"
)

type EqualNodeController struct {
	enable bool
}

func Register(eqndc *cloudagent.EqualNodeController) {
	config.InitConfigure(eqndc)
	core.Register(NewEqualNodeLister(eqndc.Enable))
}

func (eqndc *EqualNodeController) Cleanup() {}
func (eqndc *EqualNodeController) Name() string {
	return modules.EqualNodeControllerModule
}

func (eqndc *EqualNodeController) Group() string {
	return modules.EqualNodeControllerGroup
}

func (eqndc *EqualNodeController) Enable() bool {
	return eqndc.enable
}

func (eqndc *EqualNodeController) Start() {
	flag.Parse()
	go controller.StartEqualNodecontroller()
}

func NewEqualNodeLister(enable bool) *EqualNodeController {
	return &EqualNodeController{
		enable: enable,
	}
}

//func init() {
//	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
//	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
//}
