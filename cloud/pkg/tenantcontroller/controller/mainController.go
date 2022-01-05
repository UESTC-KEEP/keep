package tenant_main_controller

import (
	tenantinformer "keep/cloud/pkg/tenantcontroller/controller/tenant/informer"
)

func StartTenantController() {
	tenantinformer.StartTenantInformer()
}
