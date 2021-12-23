package trqCrdcontroller

import trqCrdinformer "keep/cloud/pkg/tenantresourcequotacontroller/controller/informer"

func StartTrqController() {
	trqCrdinformer.StartTrqInformer()
}
