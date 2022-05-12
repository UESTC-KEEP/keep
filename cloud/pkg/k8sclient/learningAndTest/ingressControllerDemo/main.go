package main

import (
	"github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/learningAndTest/ingressControllerDemo/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//获取configurati
	//创建client
	//创建informer
	//注册事件处理方法
	//启动informer
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		// 集群内部创建config
		clusterConfig, err := rest.InClusterConfig()
		if err != nil {
			return
		}
		config = clusterConfig
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
	factory := informers.NewSharedInformerFactory(clientset, 0)
	serviceInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()

	controller := pkg.NewController(clientset, serviceInformer, ingressInformer)

	stopCh := make(chan struct{})
	factory.Start(stopCh)
	// 等待同步完成再开始controller
	factory.WaitForCacheSync(stopCh)
	controller.Run(stopCh)

	clientset.NetworkingV1().Ingresses("default")
}
