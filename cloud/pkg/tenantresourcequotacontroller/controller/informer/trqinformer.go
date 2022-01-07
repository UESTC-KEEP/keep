package trqCrdinformer

import (
	tenantresourcequotav1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/tenantresourcequota/v1alpha1"
	trqcrdinformers "github.com/UESTC-KEEP/keep/cloud/pkg/client/trq/informers/externalversions"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"k8s.io/client-go/tools/cache"
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
