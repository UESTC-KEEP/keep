package tenant_onadded

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tenantv1 "keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	"keep/cloud/pkg/common/client"
	logger "keep/pkg/util/loggerv1.0.1"
	"sync"
)

func createNamespace(in_tenant *tenantv1.Tenant) {
	pod, _ := client.GetKubeClient().CoreV1().Pods("default").Get(context.Background(), "keep-eqnd-test-nginx-7dc485b56f-67kk8", metav1.GetOptions{})
	pod.ObjectMeta.Annotations["cni.projectcalico.org/podIP"] = "172.169.214.249/32"
	pod.ObjectMeta.Annotations["cni.projectcalico.org/podIPs"] = "172.169.214.249/32"
	pod.Status.PodIP = "172.169.214.249"
	pod.Status.PodIPs = []corev1.PodIP{
		{IP: "172.169.214.249"},
	}
	_, err := client.GetKubeClient().CoreV1().Pods("default").Update(context.Background(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("修改网络失败：", err)
	}

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
		go func(new_ns string) {
			ns := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					// 用户名+命名空间名+uuid
					Name: new_ns,
				},
			}
			_, err = client.GetKubeClient().CoreV1().Namespaces().Create(context.Background(), &ns, metav1.CreateOptions{})
			if err != nil {
				logger.Error("创建ns失败：", err)
				return
			}
			wg.Done()
		}(newNs)
	}
}
