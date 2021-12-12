package pods

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"keep/cloud/pkg/k8sclient/config"
	logger "keep/pkg/util/loggerv1.0.1"
)

type PodsImpl struct {
	Namespace string `json:"namespace"`
}

func (pi *PodsImpl) ListPods(namespace string) (*corev1.PodList, error) {
	podlist, err := config.Clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}
	return podlist, err
}

// GetPodInfoByPodName function to get pod from k8s by name
func (pi *PodsImpl) GetPodInfoByPodName(podName string) (*corev1.Pod, error) {
	pod, err := config.Clientset.CoreV1().Pods(metav1.NamespaceDefault).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return pod, nil
}

func NewPods() *PodsImpl {
	return &PodsImpl{}
}

//func (pi *PodsImpl) ReDeployPodToAnotherNode() {
//	pod, err := pi.GetPodInfoByPodName("clk-beep-mapper-7dc5855684-nkn5c")
//	if err != nil {
//		return
//	}
//	//newNode := "edge-zlj-ubuntu-1"
//	//pod.Spec.NodeName
//	patchData := map[string]interface{}{"metadata": map[string]map[string]string{"spec": {
//		"nodename": "edge-zlj-ubuntu-1",
//	}}}
//
//	playLoadBytes, err := json.Marshal(patchData)
//	if err !=nil{
//		logger.Error(err)
//	}
//	eviction := &policy.Eviction{
//		TypeMeta: metav1.TypeMeta{
//			APIVersion: policyGroupVersion,
//			Kind:       eutils.EvictionKind,
//		},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      pod.Name,
//			Namespace: pod.Namespace,
//		},
//		DeleteOptions: deleteOptions,
//	}
//	pod_, err := config.Clientset.PolicyV1().Evictions().Evict(context.Background(),policy.Eviction{
//		TypeMeta:      metav1.TypeMeta{},
//		ObjectMeta:    metav1.ObjectMeta{},
//		DeleteOptions: nil,
//	})
//	if err != nil {
//		logger.Error(err)
//		return
//	}
//	fmt.Println(pod_)
//}
