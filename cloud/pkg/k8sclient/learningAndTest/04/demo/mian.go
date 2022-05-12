package main

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	// configuration
	// client
	// informers
	// 注册事件处理方法
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		// 在集群内部创建
		inclusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln("cant get config")
		}
		config = inclusterConfig
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("cant get client")
	}
	informers.NewSharedInformerFactory(clientset, 0)
}
