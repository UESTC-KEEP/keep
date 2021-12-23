package trqCrdinformer

import (
	"k8s.io/client-go/tools/cache"
	tenantresourcequotav1 "keep/cloud/pkg/apis/keepedge/tenantresourcequota/v1alpha1"
	trqcrdinformers "keep/cloud/pkg/client/trq/informers/externalversions"
	"keep/cloud/pkg/common/client"
	beehiveContext "keep/pkg/util/core/context"
	logger "keep/pkg/util/loggerv1.0.1"
	"time"
)

func StartTrqInformer() {
	logger.Debug("启动 tenantresourcequota informer...")
	trqInformer := trqcrdinformers.NewSharedInformerFactory(client.GetTrqCRDClient(), 2*time.Second).Keepedge().V1alpha1().TenantResourceQuotas().Informer()
	trqInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    OnTrqAdded,
		UpdateFunc: nil,
		DeleteFunc: OnTrqDeleted,
	})
	trqInformer.Run(beehiveContext.Done())
}

func OnTrqAdded(newTrq interface{}) {
	logger.Debug("tenantresourcequota 添加: ", newTrq.(*tenantresourcequotav1.TenantResourceQuota).Name)
}

func OnTrqDeleted(delTrq interface{}) {
	logger.Debug("tenantresourcequota 删除: ", delTrq.(*tenantresourcequotav1.TenantResourceQuota).Name)
}
