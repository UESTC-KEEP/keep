package pods

import (
	"context"
	"github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/config"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodsImpl struct {
	Namespace string `json:"namespace"`
}

func (pi *PodsImpl) ListPods(namespace string) (*corev1.PodList, error) {
	podlist, err := config.Clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		// ResourceVersion设置为0 apiserver就不用查etcd 直接用缓存
		//resourceVersion 未设置	resourceVersion="0"	resourceVersion="<非零值>"
		//最新版本					任何版本				不老于给定版本
		ResourceVersion: "0",
	})

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

func (pi *PodsImpl) ReDeployPodToAnotherNode() {
	//pod, err := pi.GetPodInfoByPodName("redis-0")
	//if err != nil {
	//	return
	//}
	//pod.Spec.NodeSelector = map[string]string{"eqnd": "true"}
	//
	//_, err = config.Clientset.CoreV1().Pods(metav1.NamespaceDefault).Bind()
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
}
