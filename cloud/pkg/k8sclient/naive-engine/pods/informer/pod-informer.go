package naive_engin_pod_informer

import (
	"fmt"
	kubeinformer "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"keep/cloud/pkg/k8sclient/config"
	beehiveContext "keep/pkg/util/core/context"
	logger "keep/pkg/util/loggerv1.0.1"
)

func StartPodInformer() {
	logger.Debug("启动 PodInformer...")
	kubeinformerFactory := kubeinformer.NewSharedInformerFactory(config.Clientset, 2000000000)
	podsinformer := kubeinformerFactory.Apps().V1().Deployments().Informer()
	podsinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    OnPodsAdded,
		UpdateFunc: nil,
		DeleteFunc: nil,
	})
	// 启动informer
	//kubeinformerFactory.Start(wait.NeverStop)
	//kubeinformerFactory.WaitForCacheSync(wait.NeverStop)
	podsinformer.Run(beehiveContext.Done())
}

func OnPodsAdded(newPod interface{}) {
	fmt.Println("新pod加入：", newPod)
}
