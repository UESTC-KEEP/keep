package controller

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"keep/cloud/pkg/client/clientset/versioned"
	"keep/cloud/pkg/client/informers/externalversions"
	"keep/cloud/pkg/equalnodecontroller/config"
	"keep/cloud/pkg/equalnodecontroller/pkg/signals"
	"time"
)

func StartEqualNodecontroller() {
	// 处理信号量
	stopCh := signals.SetupSignalHandler()
	masterURL, kubeconfig := config.Config.MasterURL, config.Config.KubeConfig
	fmt.Println(masterURL, kubeconfig)
	// 处理入参
	//K8sConfig, err = clientcmd.BuildConfigFromFlags("", Config.KubeConfigFilePath)
	// 在配置文件中有ip不再写ip否则出错
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logger.Fatal("Error building kubeconfig: ", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logger.Fatal("Error building kubernetes clientset:", err.Error())
	}

	equalnodeClient, err := versioned.NewForConfig(cfg)
	if err != nil {
		logger.Fatal("Error building  clientset:", err.Error())
	}

	equalnodeInformerFactory := externalversions.NewSharedInformerFactory(equalnodeClient, time.Second*30)

	//得到controller
	controller := NewController(kubeClient, equalnodeClient,
		equalnodeInformerFactory.Keepedge().V1().EqualNodes())

	//启动informer
	go equalnodeInformerFactory.Start(stopCh)

	//controller开始处理消息
	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
