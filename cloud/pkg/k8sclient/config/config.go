package config

import (
	"github.com/wonderivan/logger"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"keep/constants"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once
var Clientset *kubernetes.Clientset
var DCI discovery.DiscoveryInterface
var K8sConfig *restclient.Config
var DD dynamic.Interface
var GR []*restmapper.APIGroupResources

type Configure struct {
	cloudagent.K8sClient
}

func InitConfigure(k *cloudagent.K8sClient) {
	once.Do(func() {
		Config = Configure{
			K8sClient: *k,
		}
		GetClient()
	})
}

func GetClient() {
	var err error
	K8sConfig, err = clientcmd.BuildConfigFromFlags("", constants.DefaultKubeConfigPath)
	if err != nil {
		logger.Error(err.Error())
	}
	Clientset, err = kubernetes.NewForConfig(K8sConfig)
	if err != nil {
		logger.Error(err.Error())
		err = nil
	}
	DCI = Clientset.Discovery()
	DD, err = dynamic.NewForConfig(K8sConfig)
	if err != nil {
		logger.Error(err.Error())
		err = nil
	}
	GR, err = restmapper.GetAPIGroupResources(Clientset.Discovery())
	if err != nil {
		logger.Error(err.Error())
		err = nil
	}
}

func Get() *Configure {
	return &Config
}
