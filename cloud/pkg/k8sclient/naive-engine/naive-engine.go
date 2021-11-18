package naive_engine

import (
	"bytes"
	"context"
	"fmt"
	"github.com/wonderivan/logger"
	"io/ioutil"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"keep/cloud/pkg/k8sclient/config"
)

func CreatePod() {
	//将配置信息赋值给deloymentClient
	deploymentClient := config.Clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	//构建deployment
	result, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	logger.Info("create Pod Name:" + result.GetObjectMeta().GetName())
	fmt.Printf("Create Pod Name : %q \n", result.GetObjectMeta().GetName())
}

func CreateConfigMap(configmap *apiv1.ConfigMap) {
	createConfigmap, err := config.Clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), configmap, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err)
	}
	fmt.Println(createConfigmap)
}

func CreatResourcesByYAML(yamlFileName, namespace string) {
	var err error
	filebytes, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		logger.Error(err)
	}

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(filebytes), 100)
	var rawObj runtime.RawExtension
	if err = decoder.Decode(&rawObj); err != nil {
		logger.Error(err)
	}
	obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		logger.Fatal(err)
	}

	unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

	if err != nil {
		logger.Fatal(err)
	}

	mapper := restmapper.NewDiscoveryRESTMapper(config.GR)
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		logger.Error(err)
	}

	var dri dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if unstructuredObj.GetNamespace() == "" {
			unstructuredObj.SetNamespace(namespace)
		}
		dri = config.DD.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
	} else {
		dri = config.DD.Resource(mapping.Resource)
	}
	obj2, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err)
		return
	}
	fmt.Printf("%s/%s created", obj2.GetKind(), obj2.GetName())

}
