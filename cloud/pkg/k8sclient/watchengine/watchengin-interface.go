package watchengine

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"time"
)

type WatcherEngineInterface interface {
	//WriteStructIntoYAML 将传入的结构体转化为指定路径下的yaml文件
	/*
		传入参数：filepath：期望写入的文件路径
				expireTime:期望文件的保存时间 传入-1代表永久
				resource:传入需要进行转化yaml的结构体
	*/
	WriteStructIntoYAML(resource interface{}, filepath string, expireTime time.Duration) error
	//ReadStructFromYAML 将yaml文件中的yaml对象转化成结构体
	ReadStructFromYAML(resource interface{}, filepath string) (unstructured.Unstructured, error)
	// InitK8sClientWatchEngine 初始化K8sClientWatchEngine
	/**/
	InitK8sClientWatchEngine()

	// CreatResourcesByYAML 根据yaml文件创建资源
	/*
		传入参数：yamlFilepath:所需创建的资源的定义yaml文件
			    namespace：创建资源所在ns
		返回值：表征创建资源是否成功 即报错信息
	*/
	CreatResourcesByYAML(yamlFilepath, namespace string) (bool, error)
	// DeleteResourceByYAML 根据yaml文件删除资源
	/*
		传入参数：yamlFilepath:所需创建的资源的定义yaml文件
				namespace:创建资源所在ns
		返回值：表征创建资源是否成功 即报错信息
	*/
	DeleteResourceByYAML(yamlFilepath, namespace string) (bool, error)
	// ApplyResourceByYAML 根据yaml文件应用资源更新
	/*
		传入参数：yamlFilepath:所需创建的资源的定义yaml文件
				namespace:创建资源所在ns
		返回值：表征更新资源是否成功 即报错信息
	*/
	ApplyResourceByYAML(yamlFilepath, namespace string) (bool, error)
}
