package pkg

import (
	"context"
	v12 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informer "k8s.io/client-go/informers/core/v1"
	netInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreLister "k8s.io/client-go/listers/core/v1"
	v1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"time"
)

const workerNum = 5

type controller struct {
	client        kubernetes.Interface
	ingressLister v1.IngressLister
	serviceLister coreLister.ServiceLister
	queue         workqueue.RateLimitingInterface
}

func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.queue.Add(key)
}

func (c *controller) AddService(obj interface{}) {
	c.enqueue(obj)
}

func (c *controller) UpdateService(oldObj interface{}, newObj interface{}) {
	//todo 比较annotation
	if reflect.DeepEqual(oldObj, newObj) {
		// 完全相同不再处理
		return
	}
	c.enqueue(newObj)
}

func (c *controller) Run(stopCh chan struct{}) {
	//调用多个worker
	for i := 0; i < workerNum; i++ {
		// 在stopch被关闭之前保持规最少5个worker在工作
		go wait.Until(c.worker, time.Minute, stopCh)
	}
	<-stopCh
}

// 不停的从workqueue中获取任务进行执行
func (c *controller) worker() {
	for c.processNextItem() {

	}
}

func (c *controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	// 处理完之后需要把这个item从队列移除
	defer c.queue.Done(item)

	key := item.(string)
	err := c.syncService(key)
	if err != nil {
		c.handlerError(key, err)
		return false
	}
	return false
}

// 调谐资源的状态函数
func (c *controller) syncService(item string) error {
	ns, name, err := cache.SplitMetaNamespaceKey(item)
	if err != nil {
		return err
	}
	// 删除的情况  不用处理
	service, err := c.serviceLister.Services(ns).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	// 新增和删除
	_, ok := service.GetAnnotations()["ingress/http"]
	// 要注意这里默认了创建ingress的是否会使得service和ingress的名字一致 所以才能用一样的
	ingress, err := c.ingressLister.Ingresses(ns).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	// ingress should不存在就需要创建ingress
	if ok && errors.IsNotFound(err) {
		// 如果键值对存在但是找不到对应的ingress 那么就是说ingress被删除了  重新创建
		ig := c.constructIngress(ns, name)
		c.client.NetworkingV1().Ingresses(ns).Create(context.TODO(), ig, v13.CreateOptions{})
	} else if !ok && ingress != nil {
		// 如果没有这个键值对但是ingress又还在 就删除对应的ingress
		c.client.NetworkingV1().Ingresses(ns).Delete(context.TODO(), name, v13.DeleteOptions{})
	}
	return nil
}

func (c *controller) handlerError(item string, err error) {
	// 如果err出现就重新入队
	// 最多重试次数
	if c.queue.NumRequeues(item) < 10 {
		c.queue.AddRateLimited(item)
		return
	}
	runtime.HandleError(err)
	// 不再记录重试次数
	c.queue.Forget(item)
}

func (c *controller) DeleteIngress(obj interface{}) {
	ingress := obj.(*v12.Ingress)
	// 获取对应ingress的service（被控制对象）
	ownerReference := v13.GetControllerOf(ingress)
	if ownerReference == nil {
		return
	}
	if ownerReference.Kind != "Service" {
		return
	}
	// 如果是service就需要重新入队进行处理
	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *controller) constructIngress(ns string, name string) *v12.Ingress {
	PathType := v12.PathTypePrefix

	return &v12.Ingress{
		TypeMeta: v13.TypeMeta{},
		ObjectMeta: v13.ObjectMeta{
			Name:                       name,
			GenerateName:               "",
			Namespace:                  ns,
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          v13.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ZZZ_DeprecatedClusterName:  "",
			ManagedFields:              nil,
		},
		Spec: v12.IngressSpec{
			IngressClassName: nil,
			DefaultBackend:   nil,
			TLS:              nil,
			Rules: []v12.IngressRule{
				{
					Host: "example.com",
					IngressRuleValue: v12.IngressRuleValue{
						HTTP: &v12.HTTPIngressRuleValue{
							Paths: []v12.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &PathType,
									Backend: v12.IngressBackend{
										Service: &v12.IngressServiceBackend{
											Name: name,
											Port: v12.ServiceBackendPort{
												Name:   "http-port",
												Number: 80,
											},
										},
										Resource: nil,
									},
								},
							},
						},
					},
				},
			},
		},
		Status: v12.IngressStatus{},
	}
}

func NewController(client kubernetes.Interface, serviceInformer informer.ServiceInformer, ingressInformer netInformer.IngressInformer) controller {
	c := controller{
		client:        client,
		ingressLister: ingressInformer.Lister(),
		serviceLister: serviceInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.AddService,
		UpdateFunc: c.UpdateService,
	})
	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.DeleteIngress,
	})
	return c
}
