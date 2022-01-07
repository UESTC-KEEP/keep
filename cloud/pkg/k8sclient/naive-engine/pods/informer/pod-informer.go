package naive_engin_pod_informer

import (
	"github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/config"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	corev1 "k8s.io/api/core/v1"
	kubeinformer "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func StartPodInformer() {
	logger.Debug("启动 PodInformer...")
	kubeinformerFactory := kubeinformer.NewSharedInformerFactory(config.Clientset, 2000000000)
	podsinformer := kubeinformerFactory.Core().V1().Pods().Informer()
	podsinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    OnPodsAdded,
		UpdateFunc: nil,
		DeleteFunc: OnPodDeleted,
	})
	// 启动informer
	//kubeinformerFactory.Start(wait.NeverStop)
	//kubeinformerFactory.WaitForCacheSync(wait.NeverStop)
	podsinformer.Run(beehiveContext.Done())
}

func OnPodsAdded(newPod interface{}) {
	logger.Debug("新pod加入：", newPod.(*corev1.Pod).Name)
	//fmt.Println( newPod)
}

func OnPodDeleted(delPod interface{}) {
	logger.Debug("Pod被删除：", delPod.(*corev1.Pod).Name)
	//fmt.Println(delPod)
}
