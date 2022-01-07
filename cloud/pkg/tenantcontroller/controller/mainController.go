package tenant_main_controller

import (
	tenantinformer "github.com/UESTC-KEEP/keep/cloud/pkg/tenantcontroller/controller/tenant/informer"
)

func StartTenantController() {
	tenantinformer.StartTenantInformer()
}
