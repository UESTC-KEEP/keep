package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TenantSpec `json:"spec"`
}

const (
	Initializing = "initializing"
	Pending      = "pending"
	Active       = "active"
	Failed       = "failed"
	Terminating  = "terminating"
)

type TenantSpec struct {
	// 创建Ｔｅｎａｎｔ资源对象时新建租户管理员用户的用户名
	Username string `json:"username"`
	// 创建Ｔｅｎａｎｔ资源对象时新建租户管理员用户的密码
	Password string `json:"password"`
	// 通过指定ｔｅｎａｎｔＩＤ复用Ｋｅｙｓｔｏｎｅ中己有的租户
	TenantID string `json:"tenant_id"`
	// Ｔｅｎａｎｔ资源对象所处的状态
	Status string `json:"status"`
	// 显示Ｔｅｎａｎｔ资源对象处于当前状态的原因
	Message string `json:"message"`
}

// TenantResourceQuotaList is a list of TenantResourceQuota resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Tenant `json:"items"`
}
