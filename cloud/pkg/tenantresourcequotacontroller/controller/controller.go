package trqCrdcontroller

import trqCrdinformer "github.com/UESTC-KEEP/keep/cloud/pkg/tenantresourcequotacontroller/controller/informer"

func StartTrqController() {
	trqCrdinformer.StartTrqInformer()
}
