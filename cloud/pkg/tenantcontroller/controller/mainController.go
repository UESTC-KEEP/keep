package tenant_main_controller

import (
	tenant_controller_queen "github.com/UESTC-KEEP/keep/cloud/pkg/tenantcontroller/controller/tenant"
)

func StartTenantController() {
	tenant_controller_queen.StartTenantController()
}
