package tenantresourcequotacontroller

import (
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/tenantresourcequotacontroller/config"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/core"
	logger "keep/pkg/util/loggerv1.0.1"
)

type tenantresourcequotacontroller struct {
	enable bool
}

func Register(trq *cloudagent.TenantResourceQuotaController) {
	config.InitConfigure(trq)
	trq_, err := newTenantResourceQuotaController(trq.Enable)
	if err != nil {
		logger.Error(err)
	}
	core.Register(trq_)
}

func (trc *tenantresourcequotacontroller) Cleanup() {}

func (trc *tenantresourcequotacontroller) Name() string {
	return modules.TenantResourceQuotaControllerModule
}

func (trc *tenantresourcequotacontroller) Group() string {
	return modules.TenantResourceQuotaControllerGroup
}

func (trc *tenantresourcequotacontroller) Enable() bool {
	return trc.enable
}

func (trc *tenantresourcequotacontroller) Start() {

}

func newTenantResourceQuotaController(enable bool) (*tenantresourcequotacontroller, error) {
	return &tenantresourcequotacontroller{
		enable: enable,
	}, nil
}
