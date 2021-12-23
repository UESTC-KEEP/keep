package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TenantResourceQuota struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TenantResourceQuotaSpec `json:"spec"`
}

type TenantResourceQuotaSpec struct {
	Name           string         `json:"name"`
	ResourceQuotas ResourceQuotas `json:"resourcequotas"`
}

// TenantResourceQuotaList is a list of TenantResourceQuota resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TenantResourceQuotaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []TenantResourceQuota `json:"items"`
}

type ResourceQuotas struct {
	Pods string `json:"pods"`
	Cpu  string `json:"cpu"`
}
