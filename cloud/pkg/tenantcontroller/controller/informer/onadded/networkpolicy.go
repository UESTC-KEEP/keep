package tenant_onadded

import (
	tenantv1 "keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
)

func createNetworkPolicy(newtenant *tenantv1.Tenant) {
	// 默认使得租户自己的ns中的所有流量互通
	//v1.NetworkPolicy{
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name:      "",
	//		Namespace: "",
	//	},
	//	Spec: v1.NetworkPolicySpec{
	//		PodSelector: metav1.LabelSelector{},
	//		Ingress: []v1.NetworkPolicyIngressRule{
	//			{
	//				Ports: nil,
	//				From: []v1.NetworkPolicyPeer{
	//					{
	//						PodSelector:       nil,
	//						NamespaceSelector: nil,
	//						IPBlock:           nil,
	//					},
	//				},
	//			},
	//		},
	//		Egress:      nil,
	//		PolicyTypes: []v1.PolicyType{v1.PolicyTypeEgress, v1.PolicyTypeIngress},
	//	},
	//}
	//client.GetKubeClient().NetworkingV1().NetworkPolicies().Create(context.Background())
	////client.GetKubeClient().NetworkingV1().NetworkPolicies()
}
