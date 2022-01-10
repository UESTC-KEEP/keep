package tenant_onadded

import (
	"context"
	tenantv1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	"github.com/UESTC-KEEP/keep/cloud/pkg/tenantcontroller/constant"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createNetworkPolicy(newtenant *tenantv1.Tenant) {
	// 默认使得租户自己的ns中的所有流量互通  拒绝其他ns的所有流量 允许ns出去到任何地方
	netpolicy := v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      newtenant.Spec.TenantID + "-policy",
			Namespace: "",
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					Ports: nil,
					From: []v1.NetworkPolicyPeer{
						{
							PodSelector: nil,
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{constant.NamespaceLableName: newtenant.Spec.TenantID},
							},
							IPBlock: nil,
						},
					},
				},
			},
			Egress:      []v1.NetworkPolicyEgressRule{},
			PolicyTypes: []v1.PolicyType{v1.PolicyTypeEgress, v1.PolicyTypeIngress},
		},
	}
	_, err := client.GetKubeClient().NetworkingV1().NetworkPolicies(corev1.NamespaceDefault).Create(context.Background(), &netpolicy, metav1.CreateOptions{})
	if err != nil {
		logger.Error("绑定租户网络策略失败：", err)
	}
}
