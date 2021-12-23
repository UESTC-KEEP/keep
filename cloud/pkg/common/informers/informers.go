package informers

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic/dynamicinformer"
	k8sinformer "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	crdinformers "keep/cloud/pkg/client/eqnd/informers/externalversions"
	"keep/cloud/pkg/common/client"
	"keep/pkg/util/loggerv1.0.1"
	"sync"
	"time"
)

type newInformer func() cache.SharedIndexInformer

type Manager interface {
	GetK8sInformerFactory() k8sinformer.SharedInformerFactory
	GetCRDInformerFactory() crdinformers.SharedInformerFactory
	GetDynamicSharedInformerFactory() dynamicinformer.DynamicSharedInformerFactory
	Start(stopCh <-chan struct{})
}

type informers struct {
	defaultResync                time.Duration
	keepClient                   kubernetes.Interface
	lock                         sync.Mutex
	informers                    map[string]cache.SharedIndexInformer
	crdSharedInformerFactory     crdinformers.SharedInformerFactory
	k8sSharedInformerFactory     k8sinformer.SharedInformerFactory
	dynamicSharedInformerFactory dynamicinformer.DynamicSharedInformerFactory
}

var globalInformers Manager
var once sync.Once

func GetInformersManager() Manager {
	once.Do(func() {
		globalInformers = &informers{
			defaultResync:                0,
			keepClient:                   client.GetKubeClient(),
			informers:                    make(map[string]cache.SharedIndexInformer),
			crdSharedInformerFactory:     crdinformers.NewSharedInformerFactory(client.GetEqndCRDClient(), 0),
			k8sSharedInformerFactory:     k8sinformer.NewSharedInformerFactory(client.GetKubeClient(), 0),
			dynamicSharedInformerFactory: dynamicinformer.NewFilteredDynamicSharedInformerFactory(client.GetDynamicClient(), 0, v1.NamespaceAll, nil),
		}
	})
	return globalInformers
}

func (ifs *informers) GetK8sInformerFactory() k8sinformer.SharedInformerFactory {
	return ifs.k8sSharedInformerFactory
}

func (ifs *informers) GetCRDInformerFactory() crdinformers.SharedInformerFactory {
	return ifs.crdSharedInformerFactory
}

func (ifs *informers) GetDynamicSharedInformerFactory() dynamicinformer.DynamicSharedInformerFactory {
	return ifs.dynamicSharedInformerFactory
}

func (ifs *informers) Start(stopCh <-chan struct{}) {
	ifs.lock.Lock()
	defer ifs.lock.Unlock()

	for name, informer := range ifs.informers {
		logger.Info("start informer ", name)
		go informer.Run(stopCh)
	}
	ifs.k8sSharedInformerFactory.Start(stopCh)
	ifs.crdSharedInformerFactory.Start(stopCh)
	ifs.dynamicSharedInformerFactory.Start(stopCh)
}

// getInformer get a informer named "name" or store a informer got by "newFunc" as key "name"
func (ifs *informers) getInformer(name string, newFunc newInformer) cache.SharedIndexInformer {
	ifs.lock.Lock()
	defer ifs.lock.Unlock()
	informer, exist := ifs.informers[name]
	if exist {
		return informer
	}
	informer = newFunc()
	ifs.informers[name] = informer
	return informer
}
