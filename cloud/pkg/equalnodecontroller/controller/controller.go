package controller

// $GOPATH/src/k8s.io/code-generator/generate-groups.sh all keep/cloud/pkg/client keep/cloud/pkg/k8sclient/crd_engin/keepedge/pkg/apis keepedge:v1
import (
	"fmt"
	keepcrdv1 "keep/cloud/pkg/apis/keepedge/v1"
	"keep/cloud/pkg/equalnodecontroller/constants"
	"keep/pkg/util/kplogger"
	"keep/pkg/util/loggerv1.0.1"
	"time"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	clientset "keep/cloud/pkg/client/clientset/versioned"
	equalnodescheme "keep/cloud/pkg/client/clientset/versioned/scheme"
	informers "keep/cloud/pkg/client/informers/externalversions/keepedge/v1"
	listers "keep/cloud/pkg/client/listers/keepedge/v1"
)

// Controller is the controller implementation for Student resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// equalnodeclientset is a clientset for our own API group
	equalnodeclientset clientset.Interface
	equalnodesLister   listers.EqualNodeLister
	equalnodesSynced   cache.InformerSynced
	workqueue          workqueue.RateLimitingInterface
	recorder           record.EventRecorder
}

// NewController returns a new student controller
func NewController(
	kubeclientset kubernetes.Interface,
	equalnodeclientset clientset.Interface,
	equalnodeInformer informers.EqualNodeInformer) *Controller {

	utilruntime.Must(equalnodescheme.AddToScheme(scheme.Scheme))
	logger.Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(kplogger.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: constants.ControllerAgentName})

	controller := &Controller{
		kubeclientset:      kubeclientset,
		equalnodeclientset: equalnodeclientset,
		equalnodesLister:   equalnodeInformer.Lister(),
		equalnodesSynced:   equalnodeInformer.Informer().HasSynced,
		workqueue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EqualNodes"),
		recorder:           recorder,
	}

	logger.Info("Setting up event handlers")
	// Set up an event handler for when Student resources change
	equalnodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueEqualNode,
		UpdateFunc: func(old, new interface{}) {
			oldEqualNode := old.(*keepcrdv1.EqualNode)
			newEqualNode := new.(*keepcrdv1.EqualNode)
			if oldEqualNode.ResourceVersion == newEqualNode.ResourceVersion {
				//版本一致，就表示没有实际更新的操作，立即返回
				return
			}
			controller.enqueueEqualNode(new)
		},
		DeleteFunc: controller.enqueueEqualNodeForDelete,
	})

	return controller
}

// Run 在此处开始controller的业务
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	logger.Info("开始controller业务，开始一次缓存数据同步")
	if ok := cache.WaitForCacheSync(stopCh, c.equalnodesSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	logger.Info("worker启动")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("worker已经启动")
	<-stopCh
	glog.Info("worker已经结束")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// 取数据处理
func (c *Controller) processNextWorkItem() bool {

	obj, shutdown := c.workqueue.Get()
	//event
	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {

			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// 在syncHandler中处理业务
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		logger.Info("Successfully synced ", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// 处理
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// 从缓存中取对象
	equalnode, err := c.equalnodesLister.EqualNodes(namespace).Get(name)
	if err != nil {
		// 如果Student对象被删除了，就会走到这里，所以应该在这里加入执行
		if errors.IsNotFound(err) {
			logger.Info("Student对象被删除，请在这里执行实际的删除业务: %s ", namespace, " ... ", name)

			return nil
		}

		runtime.HandleError(fmt.Errorf("failed to list equalnode by: %s/%s", namespace, name))

		return err
	}
	logger.Info("这里是equalnode对象的期望状态:  ", equalnode, "  ...")
	logger.Info("实际状态是从业务层面得到的，此处应该去的实际状态，与期望状态做对比，并根据差异做出响应(新增或者删除)")

	c.recorder.Event(equalnode, corev1.EventTypeNormal, constants.SuccessSynced, constants.MessageResourceSynced)
	return nil
}

// 数据先放入缓存，再入队列
func (c *Controller) enqueueEqualNode(obj interface{}) {
	var key string
	var err error
	// 将对象放入缓存
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	// 将key放入队列
	c.workqueue.AddRateLimited(key)
}

// 删除操作
func (c *Controller) enqueueEqualNodeForDelete(obj interface{}) {
	var key string
	var err error
	// 从缓存中删除指定对象
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	//再将key放入队列
	c.workqueue.AddRateLimited(key)
}
