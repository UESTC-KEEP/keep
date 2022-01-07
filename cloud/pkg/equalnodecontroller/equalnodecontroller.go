package equalnodecontroller

import (
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/informers"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/modules"
	"github.com/UESTC-KEEP/keep/cloud/pkg/equalnodecontroller/config"
	"github.com/UESTC-KEEP/keep/cloud/pkg/equalnodecontroller/controller"
	"github.com/UESTC-KEEP/keep/cloud/pkg/equalnodecontroller/controller/lister"
	cloudagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"github.com/UESTC-KEEP/keep/pkg/util/core"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	flag "github.com/spf13/pflag"
)

type EqualNodeController struct {
	equalnodecontroller *controller.EqualNodeController
	enable              bool
}

func Register(eqndc *cloudagent.EqualNodeController) {
	config.InitConfigure(eqndc)
	core.Register(newEqualNodeLister(eqndc.Enable))
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
	if err := eqndc.equalnodecontroller.Start(); err != nil {
		logger.Fatal("启动 equalnodecontroller crd 失败...", err)
	}
	logger.Debug(lister.GetAllEqnd())
	//go controller.StartEqualNodecontroller()
}

func newEqualNodeLister(enable bool) *EqualNodeController {
	eqndctl, err := controller.NewEqualNodeController(informers.GetInformersManager().GetCRDInformerFactory())
	if err != nil {
		logger.Fatal("创建equalnodecontroller 失败：", err)
	}
	return &EqualNodeController{
		equalnodecontroller: eqndctl,
		enable:              enable,
	}
}
