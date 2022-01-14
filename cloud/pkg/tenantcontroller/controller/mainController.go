package tenant_main_controller

import (
	tenant_controller_queen "github.com/UESTC-KEEP/keep/cloud/pkg/tenantcontroller/controller/tenant/informer"
)

func StartTenantController() {
	tenant_controller_queen.StartTenantController()
}
