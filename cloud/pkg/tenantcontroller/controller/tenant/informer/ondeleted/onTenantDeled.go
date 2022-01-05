package tenant_ondeleted

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tenantv1 "keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	"keep/cloud/pkg/common/client"
	logger "keep/pkg/util/loggerv1.0.1"
)

func DeleteTenant(deltenant *tenantv1.Tenant) {
	// 删除对应命名空间
	delTenantNamespaces(deltenant)
}

func delTenantNamespaces(deltenant *tenantv1.Tenant) {
	fmt.Println(deltenant.Spec.Namespaces)
	for _, ns := range deltenant.Spec.Namespaces {
		logger.Trace("删除命名空间：", ns)
		err := client.GetKubeClient().CoreV1().Namespaces().Delete(context.Background(), ns, metav1.DeleteOptions{})
		if err != nil {
			logger.Error(err)
		}
	}
}
