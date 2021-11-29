package equalnodecontroller

import (
	"fmt"
	"github.com/golang/glog"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"keep/cloud/pkg/client/clientset/versioned"
	"keep/cloud/pkg/client/informers/externalversions"
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/equalnodecontroller/config"
	"keep/cloud/pkg/equalnodecontroller/pkg/signals"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/core"
	"time"
)

type EqualNodeController struct {
	enable bool
}

func Register(eqndc *cloudagent.EqualNodeController) {
	config.InitConfigure(eqndc)
	core.Register(NewEqualNodeLister(eqndc.Enable))
}

func (eqndc *EqualNodeController) Cleanup() {}
func (eqndc *EqualNodeController) Name() string {
	return modules.EqualNodeControllerModule
}

func (eqndc *EqualNodeController) Group() string {
	return modules.EqualNodeControllerGroup
}

func (eqndc *EqualNodeController) Enable() bool {
	return eqndc.enable
}

func (eqndc *EqualNodeController) Start() {
	flag.Parse()
	go StartEqualNodecontroller()
}

func NewEqualNodeLister(enable bool) *EqualNodeController {
	return &EqualNodeController{
		enable: enable,
	}
}

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
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	equalnodeClient, err := versioned.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building  clientset: %s", err.Error())
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

//func init() {
//	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
//	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
//}
