package naive_engine

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Node interface {
	// GetSecret 获取secret信息
	/*
		传入参数：secretName:所需查询的secret名字
		        namespace:所需查询的secret所属命名空间
	*/
	GetSecret(secretName, namespace string) (unstructured.Unstructured, error)
	// ListPods 得到某个命名空间下所有的Pod
	/*
		传入参数：namespace：所需查询的命名空间  如果是""则使用constant.DefaultNameSpace
		返回值：得到的list   有错返回 无错nil
	*/
	ListPods(namespace string) (unstructured.UnstructuredList, error)
	// GetPodInfoByPodName 根据podname在指定的ns中查询pod的详细内容
	/*
		传入参数：podName:查询的pod的名字 考虑到一个pod可能多副本返回list
	*/
	GetPodInfoByPodName(podName string) (unstructured.UnstructuredList, error)
}
