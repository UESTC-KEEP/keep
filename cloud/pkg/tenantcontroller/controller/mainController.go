package tenant_main_controller

import tenantinformer "keep/cloud/pkg/tenantcontroller/controller/informer"

func StartTenantController() {
	tenantinformer.StartTenantInformer()
}
