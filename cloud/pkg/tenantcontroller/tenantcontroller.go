package tenant_controller

import (
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/tenantcontroller/config"
	tenant_main_controller "keep/cloud/pkg/tenantcontroller/controller"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/core"
	"keep/pkg/util/loggerv1.0.1"
	"os"
)

type TenantController struct {
	enable bool
}

func Register(tc *cloudagent.TenantController) {
	config.InitConfigure(tc)
	tenantcontroller, err := NewTenantController(tc.Enable)
	if err != nil {
		logger.Error("初始化promserver失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(tenantcontroller)
}

func (tc *TenantController) Cleanup() {}

func (tc *TenantController) Start() {
	tenant_main_controller.StartTenantController()
}

func (tc *TenantController) Name() string {
	return modules.TenantControllerModule
}
func (tc *TenantController) Group() string {
	return modules.TenantControllerGroup
}

func (tc *TenantController) Enable() bool {
	return tc.enable
}

func NewTenantController(enable bool) (*TenantController, error) {
	return &TenantController{
		enable: enable,
	}, nil
}
