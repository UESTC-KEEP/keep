package tenantresourcequotacontroller

import (
	//crdinformers "keep/cloud/pkg/client/trq/informers/externalversions"
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/tenantresourcequotacontroller/config"
	trqCrdcontroller "keep/cloud/pkg/tenantresourcequotacontroller/controller"

	//"keep/cloud/pkg/tenantresourcequotacontroller/manager"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/core"
	logger "keep/pkg/util/loggerv1.0.1"
)

type tenantresourcequotacontroller struct {
	enable bool
}

func Register(trqc *cloudagent.TenantResourceQuotaController) {
	config.InitConfigure(trqc)
	trq_, err := newTenantResourceQuotaController(trqc.Enable)
	if err != nil {
		logger.Error(err)
	}
	core.Register(trq_)
}

func (trqc *tenantresourcequotacontroller) Cleanup() {}

func (trqc *tenantresourcequotacontroller) Name() string {
	return modules.TenantResourceQuotaControllerModule
}

func (trqc *tenantresourcequotacontroller) Group() string {
	return modules.TenantResourceQuotaControllerGroup
}

func (trqc *tenantresourcequotacontroller) Enable() bool {
	return trqc.enable
}

func (trqc *tenantresourcequotacontroller) Start() {
	go trqCrdcontroller.StartTrqController()
}

func newTenantResourceQuotaController(enable bool) (*tenantresourcequotacontroller, error) {
	return &tenantresourcequotacontroller{
		enable: enable,
	}, nil
}
