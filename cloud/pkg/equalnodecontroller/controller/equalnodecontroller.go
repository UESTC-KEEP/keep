package controller

import (
	"github.com/golang/glog"
	"github.com/wonderivan/logger"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	crdv1 "keep/cloud/pkg/apis/keepedge/v1"
	"keep/cloud/pkg/client/clientset/versioned"
	"keep/cloud/pkg/client/informers/externalversions"
	crdinformers "keep/cloud/pkg/client/informers/externalversions"
	"keep/cloud/pkg/common/client"
	"keep/cloud/pkg/equalnodecontroller/config"
	"keep/cloud/pkg/equalnodecontroller/manager"
	"keep/cloud/pkg/equalnodecontroller/pkg/signals"
	beehiveContext "keep/pkg/util/core/context"
	"reflect"
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
	equalnodeManager, err := manager.NewEqualNodeManager(crdInformerFactory.Keepedge().V1().EqualNodes().Informer())
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
	logger.Error("开始监听......................")
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

func (eqndctl *EqualNodeController) equalNodeAdded(eqnd *crdv1.EqualNode) {
	eqndctl.equalnodeManager.EqualNode.Store(eqnd.Name, eqnd)
	logger.Info("----------- crd增加:  ")
	logger.Info("----------- crd增加:  ", eqnd)
	logger.Error(eqnd.Spec.Name, "      ", eqnd.Spec.Eqnd)
}

func (eqndctl *EqualNodeController) equalNodeDeleted(eqnd *crdv1.EqualNode) {
	eqndctl.equalnodeManager.EqualNode.Delete(eqnd.Name)
	logger.Info("----------- crd删除:  ")
	logger.Info("----------- crd删除:  ", eqnd)
}

func (eqndctl *EqualNodeController) equalNodeUpated(eqnd *crdv1.EqualNode) {
	value, ok := eqndctl.equalnodeManager.EqualNode.Load(eqnd.Name)
	eqndctl.equalnodeManager.EqualNode.Store(eqnd.Name, eqnd)
	if ok {
		cacheEqualNode := value.(*crdv1.EqualNode)
		if isEqualNodeUpdated(cacheEqualNode, eqnd) {
			logger.Info("----------- crd更新:  ")
			logger.Info("----------- crd更新:  ", eqnd)
		}
	}

}

// isDeviceUpdated 检查eqnd crd是否更新
func isEqualNodeUpdated(oldeqnd *crdv1.EqualNode, neweqnd *crdv1.EqualNode) bool {
	// does not care fields
	oldeqnd.ObjectMeta.ResourceVersion = neweqnd.ObjectMeta.ResourceVersion
	oldeqnd.ObjectMeta.Generation = neweqnd.ObjectMeta.Generation
	// return true if ObjectMeta or Spec or Status changed, else false
	return !reflect.DeepEqual(oldeqnd.ObjectMeta, neweqnd.ObjectMeta) || !reflect.DeepEqual(oldeqnd.Spec, neweqnd.Spec)
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
