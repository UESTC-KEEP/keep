package config

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"keep/cloud/pkg/common/client"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/loggerv1.0.1"
	"sync"
)

var Config Configure
var once sync.Once
var Clientset *kubernetes.Clientset
var DiscoveryClient discovery.DiscoveryInterface
var K8sConfig *restclient.Config
var DynamicClient dynamic.Interface
var ApiGroupResources []*restmapper.APIGroupResources

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
	// 在配置文件中有ip不再写ip否则出错
	K8sConfig, err = clientcmd.BuildConfigFromFlags("", Config.KubeAPIConfig.KubeConfig)
	if err != nil {
		logger.Error(err.Error())
	}
	//Clientset, err = kubernetes.NewForConfig(K8sConfig)
	Clientset = client.GetKubeClient().(*kubernetes.Clientset)
	if err != nil {
		logger.Error(err.Error())
		err = nil
	}
	DiscoveryClient = Clientset.Discovery()
	DynamicClient, err = dynamic.NewForConfig(K8sConfig)
	if err != nil {
		logger.Error(err.Error())
		err = nil
	}
	ApiGroupResources, err = restmapper.GetAPIGroupResources(Clientset.Discovery())
	if err != nil {
		logger.Error(err.Error())
		err = nil
	}
}
func GetGVRdyClient(gvk *schema.GroupVersionKind, namespace string) (dr dynamic.ResourceInterface, err error) {
	//et GVK GVR mapper
	GetClient()
	mapperGVRGVK := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(DiscoveryClient))
	resourceMapper, err := mapperGVRGVK.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		logger.Error(err.Error())
	}
	if resourceMapper.Scope.Name() == meta.RESTScopeNameNamespace {
		dr = DynamicClient.Resource(resourceMapper.Resource).Namespace(namespace)
	} else {
		dr = DynamicClient.Resource(resourceMapper.Resource)
	}
	return dr, err
}

func Get() *Configure {
	return &Config
}
