package crd_engin

//
//import (
//	"bytes"
//	"context"
//	"github.com/wonderivan/logger"
//	"io/ioutil"
//	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
//	"k8s.io/apimachinery/pkg/runtime"
//	"k8s.io/apimachinery/pkg/runtime/schema"
//	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
//	"k8s.io/client-go/dynamic"
//	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
//	studentv1 "keep/cloud/pkg/crd/client/clientset/versioned/typed/bolingcavalry/v1"
//	"keep/cloud/pkg/k8sclient/config"
//	naive_engine "keep/cloud/pkg/k8sclient/naive-engine"
//	"keep/constants"
//	"keep/pkg/util/kelogger"
//)
//
//type CrdEngineImpl struct{}
//
//func NewCrdEngineImpl()*CrdEngineImpl{
//	return &CrdEngineImpl{}
//}
//
//func (cei *CrdEngineImpl)CreateCrd(Dir string)error{
//	var err error
//	if err != nil {
//		kelogger.Error(err)
//	}
//	Clientset, err := studentv1.NewForConfig(config.K8sConfig)
//	naive_engine.CreatResourcesByYAML(constants.DefaultCrdsDir+"/"+"stuCrd.yaml",constants.DefaultNameSpace)
//
//	naive_engine.CreatResourcesByYAML(constants.DefaultCrdsDir+"/"+"new-student.yaml",constants.DefaultNameSpace)
//	//Clientset.Students("default").Delete(context.TODO(),"new-student",metav1.DeleteOptions{})
//	return err
//}
//
//func CreateCrdByYMAL(yamlFileName string){
//	//dynamicClient,err := dynamic.NewForConfig(config.K8sConfig)
//	//if err != nil {
//	//	kelogger.Error(err)
//	//}
//	//gvr := schema.GroupVersionResource{Version: "v1",Resource: "Student"}
//	//filebytes, err := ioutil.ReadFile(yamlFileName)
//	//if err != nil {
//	//	logger.Error(err)
//	//}
//	//
//	//decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(filebytes), 100)
//	//var rawObj runtime.RawExtension
//	//if err = decoder.Decode(&rawObj); err != nil {
//	//	logger.Error(err)
//	//}
//	//obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
//	//unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
//	//if err != nil {
//	//	logger.Fatal(err)
//	//}
//	//
//	//unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
//	//
//	//dynamicClient.Resource(gvr).Namespace("default").Create(context.TODO(),)
//}
