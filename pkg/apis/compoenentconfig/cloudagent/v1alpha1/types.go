package v1alpha1

type K8sClient struct {
	Kubeconfig string `json:"kubeconfig"`
	Kubemaster string `json:"kubemaster"`

}