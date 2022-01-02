package tenant_informer

import (
	"k8s.io/client-go/tools/cache"
	tenantv1 "keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	tenantInformers "keep/cloud/pkg/client/tenant/informers/externalversions"
	"keep/cloud/pkg/common/client"
	tenant_onadded "keep/cloud/pkg/tenantcontroller/controller/informer/onadded"
	beehiveContext "keep/pkg/util/core/context"
	logger "keep/pkg/util/loggerv1.0.1"
	"time"
)

func StartTenantInformer() {
	logger.Debug("启动tenant crd informer.....")
	tenantInformer := tenantInformers.NewSharedInformerFactory(client.GetTenantClient(), 2*time.Second).Keepedge().V1alpha1().Tenants().Informer()
	tenantInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    OnTenantAdded,
		UpdateFunc: nil,
		DeleteFunc: nil,
	})
	tenantInformer.Run(beehiveContext.Done())
}

func OnTenantAdded(newTenant interface{}) {
	newtenant := newTenant.(*tenantv1.Tenant)
	logger.Debug("新租户加入：", newtenant.Spec.Username)
	// 更新租户状态
	tenant_onadded.AddedTenant(newtenant)
}
