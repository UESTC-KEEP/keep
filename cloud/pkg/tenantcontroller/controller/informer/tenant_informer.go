package tenant_informer

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	tenantv1 "keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	tenantInformers "keep/cloud/pkg/client/tenant/informers/externalversions"
	"keep/cloud/pkg/common/client"
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
	newtenant.Spec.TenantID = "12345678901"
	newtenant.Spec.Message = "能走到这一步真的是 艰难苦恨繁霜鬓..."
	_, err := client.GetTenantClient().KeepedgeV1alpha1().Tenants(corev1.NamespaceDefault).Update(context.Background(), newtenant, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(err)
	}
}
