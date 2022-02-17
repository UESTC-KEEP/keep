package tenant_informer

import (
	"fmt"
	"github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	tenantv1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/tenant/v1alpha1"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	tenant_onadded "github.com/UESTC-KEEP/keep/cloud/pkg/tenantcontroller/controller/tenant/informer/onadded"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

/*
Copyright 2017 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Controller demonstrates how to implement a controller with client-go.
type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
}

// NewController creates a new Controller.
func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	err := c.syncToStdout(key.(string))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)
	return true
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the pod to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *Controller) syncToStdout(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		logger.Error("Fetching object with key ", key, " from store failed with ", err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a tenant, so that we will see a delete for one pod
		logger.Debug("Tenant ", key, " 被删除...")
	} else {
		// Note that you also have to check the uid if you have a local controlled resource, which
		// is dependent on the actual instance, to detect that a teant was recreated with the same name
		logger.Debug("同步/新增/更新 for Tenant ", obj.(*v1alpha1.Tenant).GetName())
	}
	return nil
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		logger.Info("Error syncing tenant ", key, ": %v", err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	logger.Info("Dropping pod ", key, " out of the queue: %v", err)
}

// Run begins watching and syncing.
func (c *Controller) Run(workers int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()
	logger.Info("Starting Tenant controller")

	go c.informer.Run(beehiveContext.Done())

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	logger.Info("Stopping Tenant controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

// 创建工作队列
var queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

func StartTenantController() {
	// 创建crd的操作客户端
	clientset := client.GetTenantClient()
	// 创建租户watcher
	tenantListWatcher := cache.NewListWatchFromClient(clientset.KeepedgeV1alpha1().RESTClient(), "tenants", v1.NamespaceDefault, fields.Everything())
	synced := cache.WaitForCacheSync(beehiveContext.Done())
	if !synced {
		logger.Error("租户控制器同步失败...")
		beehiveContext.Cancel()
	}
	logger.Debug("租户同步成功...")
	// 在infomer的帮助下把workqueue和cache绑定 确保cache更新的是否都会有key添加到workqueue
	// 需要注意的是我们最终会在workqueue中独立对象，可能会有个比触发动作更新版本的tenant
	indexer, informer := cache.NewIndexerInformer(tenantListWatcher, &v1alpha1.Tenant{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc:    OnTenantAdded,
		UpdateFunc: OnTenantUpdate,
		DeleteFunc: OnTenantDeleted,
	}, cache.Indexers{})
	controller := NewController(queue, indexer, informer)

	// 启动控制器
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// Wait forever
	select {}
}

func OnTenantAdded(newTenant interface{}) {
	if newTenant.(*v1alpha1.Tenant).Spec.Status != tenantv1.Initializing {
		logger.Debug("租户：", newTenant.(*v1alpha1.Tenant).Spec.Username, " 已被处理...")
		return
	}
	key, err := cache.MetaNamespaceKeyFunc(newTenant)
	if err == nil {
		queue.Add(key)
	}
	newtenant := newTenant.(*tenantv1.Tenant)
	logger.Debug("新租户加入：", newtenant.Spec.Username)
	// 更新租户状态
	tenant_onadded.AddedTenant(newtenant)
}

func OnTenantUpdate(old interface{}, new interface{}) {
	logger.Error(new.(*v1alpha1.Tenant).ResourceVersion, new.(*v1alpha1.Tenant).ResourceVersion)
	if new.(*v1alpha1.Tenant).ResourceVersion == old.(*v1alpha1.Tenant).ResourceVersion {
		logger.Error("原来的啥子都不干...")
		return
	}
	key, err := cache.MetaNamespaceKeyFunc(new)
	if err == nil {
		queue.Add(key)
	}
}

func OnTenantDeleted(obj interface{}) {
	// IndexerInformer uses a delta queue, therefore for deletes we have to use this
	// key function.
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err == nil {
		queue.Add(key)
	}
	deltenant := obj.(*tenantv1.Tenant)
	logger.Warn("删除租户：", deltenant.Spec.Username)

	//tenant_ondeleted.DeleteTenant(deltenant)
}
