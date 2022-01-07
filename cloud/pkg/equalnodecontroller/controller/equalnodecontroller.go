package controller

import (
	crdv1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/equalnode/v1alpha1"
	"github.com/UESTC-KEEP/keep/cloud/pkg/client/eqnd/clientset/versioned"
	crdinformers "github.com/UESTC-KEEP/keep/cloud/pkg/client/eqnd/informers/externalversions"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	"github.com/UESTC-KEEP/keep/cloud/pkg/equalnodecontroller/config"
	"github.com/UESTC-KEEP/keep/cloud/pkg/equalnodecontroller/manager"
	"github.com/UESTC-KEEP/keep/cloud/pkg/equalnodecontroller/pkg/signals"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

type EqualNodeController struct {
	kubeClient       kubernetes.Interface
	equalnodeManager *manager.EqualNodeManager
}

func (eqndctl *EqualNodeController) Start() error {
	logger.Info("开始启动 EqualNodeController ...")
	go eqndctl.TestController()
	return nil
}

func NewEqualNodeController(crdInformerFactory crdinformers.SharedInformerFactory) (*EqualNodeController, error) {
	//var workqueue workqueue.RateLimitingInterface
	//workqueue.Done()
	equalnodeManager, err := manager.NewEqualNodeManager(crdInformerFactory.Keepedge().V1alpha1().EqualNodes().Informer())
	if err != nil {
		logger.Warn("创建equalnode manager警告：", err)
		return nil, err
	}
	eqndctl := &EqualNodeController{
		kubeClient:       client.GetKubeClient(),
		equalnodeManager: equalnodeManager,
	}
	return eqndctl, nil
}

func (eqndctl *EqualNodeController) TestController() {
	logger.Info("Controller 开始监听.........")
	for {
		select {
		case <-beehiveContext.Done():
			logger.Warn("stop TestController")
			return
		case e := <-eqndctl.equalnodeManager.Events():
			equalnode, ok := e.Object.(*crdv1.EqualNode)
			if !ok {
				logger.Warn("Object type: %T unsupported:  ", equalnode)
				continue
			}
			switch e.Type {
			case watch.Added:
				eqndctl.equalNodeAdded(equalnode)
			case watch.Deleted:
				eqndctl.equalNodeDeleted(equalnode)
			case watch.Modified:
				eqndctl.equalNodeUpated(equalnode)
			default:
				logger.Warn("crd eqnd 事件类型不支持：", e.Type)
			}
		}
	}
}

func StartEqualNodecontroller() {
	// 处理信号量
	stopCh := signals.SetupSignalHandler()
	_, kubeconfig := config.Config.MasterURL, config.Config.KubeConfig
	// 在配置文件中有ip不再写ip否则出错
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logger.Fatal("Error building kubeconfig: ", err.Error())
	}

	kubeClient := client.GetKubeClient()
	if err != nil {
		logger.Fatal("Error building kubernetes clientset:", err.Error())
	}

	equalnodeClient, err := versioned.NewForConfig(cfg)
	if err != nil {
		logger.Fatal("Error building  clientset:", err.Error())
	}

	equalnodeInformerFactory := crdinformers.NewSharedInformerFactory(equalnodeClient, time.Second*30)

	//得到controller
	controller := NewController(kubeClient, equalnodeClient,
		equalnodeInformerFactory.Keepedge().V1alpha1().EqualNodes())

	//启动informer
	go equalnodeInformerFactory.Start(stopCh)

	//controller开始处理消息
	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
