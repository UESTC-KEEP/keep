package tenant_informer

import (
	tenantInformers "github.com/UESTC-KEEP/keep/cloud/pkg/client/tenant/informers/externalversions"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"k8s.io/client-go/tools/cache"
	"time"
)

func StartTenantInformer() {
	logger.Debug("启动tenant crd informer.....")
	informerFactory := tenantInformers.NewSharedInformerFactory(client.GetTenantClient(), time.Second*30)
	tenantInformer := informerFactory.Keepedge().V1alpha1().Tenants().Informer()
	tenantInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    OnTenantAdded,
		UpdateFunc: nil,
		DeleteFunc: OnTenantDeleted,
	})
	informerFactory.WaitForCacheSync(beehiveContext.Done())
	tenantInformer.Run(beehiveContext.Done())
}

//func OnTenantAdded(newTenant interface{}) {
//	newtenant := newTenant.(*tenantv1.Tenant)
//	logger.Debug("新租户加入：", newtenant.Spec.Username)
//	// 更新租户状态
//	tenant_onadded.AddedTenant(newtenant)
//}
//
//func OnTenantDeleted(delTenant interface{}) {
//	deltenant := delTenant.(*tenantv1.Tenant)
//	logger.Warn("删除租户：", deltenant.Spec.Username)
//	tenant_ondeleted.DeleteTenant(deltenant)
//}
