package tenant_onadded

import (
	"context"
	"errors"
	tenantv1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	uuid "github.com/satori/go.uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AddedTenant(newtenant *tenantv1.Tenant) {
	// TODO 存入数据库
	newStatus, err := checkTenantParameters(newtenant)
	if err != nil {
		logger.Error(err)
		UpdateTenantStatus(newtenant, newStatus)
		return
	}
	// 生成UUId
	newtenant.Spec.TenantID = (uuid.NewV4()).String()
	_, err = client.GetTenantClient().KeepedgeV1alpha1().Tenants(corev1.NamespaceDefault).Update(context.Background(), newtenant, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新租户id失败：", err)
	}

	//newNs := corev1.Create
	// 创建与租户唯一绑定的ns  设置为username '-' uuid
	createNamespace(newtenant)
	// 配置网络隔离策略  默认不与其他namespace通信 自己的ns的所有流量可以出去
	createNetworkPolicy(newtenant)
}

// 校验租户参数，不通过设置状态
func checkTenantParameters(newtenant *tenantv1.Tenant) (string, error) {
	if newtenant.Spec.Status != "" {
		return tenantv1.Failed, errors.New("you should preset tenant status,please set it as \"\"")
	}
	return tenantv1.Pending, nil
}

func UpdateTenantStatus(in_tenant *tenantv1.Tenant, status string) {
	tenant, err := client.GetTenantClient().KeepedgeV1alpha1().Tenants(corev1.NamespaceDefault).Get(context.Background(), in_tenant.Name, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取租户失败：", err)
		return
	}
	tenant.Spec.Status = status
	_, err = client.GetTenantClient().KeepedgeV1alpha1().Tenants(corev1.NamespaceDefault).Update(context.Background(), tenant, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新租户状态失败：", err)
	}
}
