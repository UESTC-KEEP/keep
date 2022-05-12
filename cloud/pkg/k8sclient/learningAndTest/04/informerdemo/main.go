package main

import (
	"fmt"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Pods().Informer()

	//添加队列
	ratelimiter := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "controller")

	//注册事件处理
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("add")
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				fmt.Println("CANT GET KEY")
			}
			ratelimiter.Add(key)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("update")
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err != nil {
				fmt.Println("CANT GET KEY")
			}
			ratelimiter.Add(key)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("add")
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				fmt.Println("CANT GET KEY")
			}
			ratelimiter.Add(key)
		},
	})
	stopch := make(chan struct{})
	// 启动
	factory.Start(stopch)
	factory.WaitForCacheSync(stopch)
	<-stopch
}
