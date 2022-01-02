package tenant_onadded

import (
	"context"
	uuid "github.com/satori/go.uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tenantv1 "keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	"keep/cloud/pkg/common/client"
	logger "keep/pkg/util/loggerv1.0.1"
)

func UpdateTenantStatus(tenant *tenantv1.Tenant, status string) {
	tenant.Spec.Status = status
	_, err := client.GetTenantClient().KeepedgeV1alpha1().Tenants(corev1.NamespaceDefault).Update(context.Background(), tenant, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(err)
	}
}

func AddedTenant(newtenant *tenantv1.Tenant) {
	UpdateTenantStatus(newtenant, tenantv1.Initializing)
	// 生成UUId
	newtenant.Spec.TenantID = (uuid.NewV4()).String()
	// 创建与租户唯一绑定的ns  设置为username '-' uuid
	// TODO - 避免ns和租户的一一绑定   允许租户创建多个ns  将不同的应用拆解到不同的ns中  与ns解耦  避免变更一个组件导致雪崩
	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: newtenant.Spec.Username + newtenant.Spec.TenantID,
		},
	}
	_, err := client.GetKubeClient().CoreV1().Namespaces().Create(context.Background(), &ns, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err)
	}
	//newNs := corev1.Create

	// TODO 存入数据库

}
