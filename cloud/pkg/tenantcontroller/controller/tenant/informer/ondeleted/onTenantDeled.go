package tenant_ondeleted

import (
	"context"
	"fmt"
	tenantv1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
