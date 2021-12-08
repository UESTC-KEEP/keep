package naive_engine

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"keep/cloud/pkg/k8sclient/config"
	"keep/pkg/util/loggerv1.0.1"
	"os"
)

type NaiveEngineImpl struct {
}

func (nei *NaiveEngineImpl) CreatePod() {
	//将配置信息赋值给deloymentClient
	deploymentClient := config.Clientset.AppsV1().Deployments(corev1.NamespaceDefault)
	//构建deployment
	result, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	logger.Info("create Pod Name:" + result.GetObjectMeta().GetName())
	fmt.Printf("Create Pod Name : %q \n", result.GetObjectMeta().GetName())
}

func (nei *NaiveEngineImpl) CreateConfigMap(configmap *corev1.ConfigMap) {
	createConfigmap, err := config.Clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), configmap, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err)
	}
	fmt.Println(createConfigmap)
}

// CreatResourcesByYAML 使用yaml文件创建资源
func (nei *NaiveEngineImpl) CreatResourcesByYAML(yamlFileName, namespace string) error {
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

	mapper := restmapper.NewDiscoveryRESTMapper(config.ApiGroupResources)
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		logger.Error(err)
	}

	var dri dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if unstructuredObj.GetNamespace() == "" {
			unstructuredObj.SetNamespace(namespace)
		}
		dri = config.DynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
	} else {
		dri = config.DynamicClient.Resource(mapping.Resource)
	}
	obj2, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err)
	} else {
		logger.Debug("%s/%s created", obj2.GetKind(), obj2.GetName())
	}
	return err
}

// DeleteResourceByYAML 使用yaml文件删除资源
func (nei *NaiveEngineImpl) DeleteResourceByYAML(filename string, namespace string) error {
	f, err := os.Open(filename)

	if err != nil {
		logger.Error(err)
		return err
	}
	d := yamlutil.NewYAMLOrJSONDecoder(f, 4096)

	restMapperRes, err := restmapper.GetAPIGroupResources(config.DiscoveryClient)
	if err != nil {
		logger.Error(err)
		return err
	}

	restMapper := restmapper.NewDiscoveryRESTMapper(restMapperRes)

	for {
		ext := runtime.RawExtension{}

		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			logrus.Fatal(err)
		}

		// runtime.Object
		obj, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
		if err != nil {
			logger.Error(err)
			return err
		}

		mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		// fmt.Printf("mapping:%+v\n", mapping)
		if err != nil {
			logger.Error(err)
			return err
		}

		// runtime.Object转换为unstructed
		unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			logger.Error(err)
			return err
		}
		// fmt.Printf("unstructuredObj: %+v", unstructuredObj)

		var unstruct unstructured.Unstructured

		unstruct.Object = unstructuredObj

		tmpMetadata := unstructuredObj["metadata"].(map[string]interface{})
		tmpName := tmpMetadata["name"].(string)
		tmpKind := unstructuredObj["kind"].(string)
		logger.Info("删除资源名：: " + tmpName + ", 资源种类: " + tmpKind + ", 所属命名空间: " + namespace)

		if namespace == "" {
			err := config.DynamicClient.Resource(mapping.Resource).Delete(context.TODO(), tmpName, metav1.DeleteOptions{})
			if err != nil {
				logger.Error(err)
				return err
			}
		} else {
			err := config.DynamicClient.Resource(mapping.Resource).Namespace(namespace).Delete(context.TODO(), tmpName, metav1.DeleteOptions{})
			if err != nil {
				logger.Error(err)
				return err
			}
		}

	}

	return nil
}

// GetSecret 按名字查询secret
func (nei *NaiveEngineImpl) GetSecret(secretName, namespace string) (*corev1.Secret, error) {
	secretlist, err := config.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secretlist, err
}

func (nei *NaiveEngineImpl) ListPods(namespace string) (*corev1.PodList, error) {
	podlist, err := config.Clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}
	return podlist, err
}

// function to get pod from k8s by name

func (nei *NaiveEngineImpl) GetPodInfoByPodName(podName string) (*corev1.Pod, error) {
	pod, err := config.Clientset.CoreV1().Pods("A").Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return pod, nil
}

func NewNaiveEngine() *NaiveEngineImpl {
	return &NaiveEngineImpl{}
}

func TestFunctions() {
	//fmt.Println(NewNaiveEngine().ListPods("default"))
}
