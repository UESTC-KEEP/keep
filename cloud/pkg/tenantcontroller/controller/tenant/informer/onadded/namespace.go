package tenant_onadded

import (
	"context"
	tenantv1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	"github.com/UESTC-KEEP/keep/cloud/pkg/tenantcontroller/constant"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"
)

func createNamespace(in_tenant *tenantv1.Tenant) {
	newtenant, err := client.GetTenantClient().KeepedgeV1alpha1().Tenants(corev1.NamespaceDefault).Get(context.Background(), in_tenant.Name, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取租户失败：", err)
		return
	}
	// TODO - 避免ns和租户的一一绑定   允许租户创建多个ns  将不同的应用拆解到不同的ns中  与ns解耦  避免变更一个组件导致雪崩
	var wg sync.WaitGroup
	// 把tenant的ns名字改成系统真实的
	for i := 0; i < len(newtenant.Spec.Namespaces); i++ {
		newtenant.Spec.Namespaces[i] = newtenant.Spec.Username + "-" + newtenant.Spec.Namespaces[i] + "-" + newtenant.Spec.TenantID
	}
	newtenant, err = client.GetTenantClient().KeepedgeV1alpha1().Tenants(corev1.NamespaceDefault).Update(context.Background(), newtenant, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新租户ns失败：", err)
		return
	}
	// 创建ns
	for _, newNs := range newtenant.Spec.Namespaces {
		wg.Add(1)
		go func(new_ns string, tenant *tenantv1.Tenant) {
			ns := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					// 用户名+命名空间名+uuid  同时打上tenantid =
					Name:   new_ns,
					Labels: map[string]string{constant.NamespaceLableName: tenant.Spec.TenantID},
				},
			}
			_, err = client.GetKubeClient().CoreV1().Namespaces().Create(context.Background(), &ns, metav1.CreateOptions{})
			if err != nil {
				logger.Error("创建ns失败：", err)
				return
			}
			wg.Done()
		}(newNs, newtenant)
	}
	wg.Wait()
}
