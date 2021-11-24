// Package watchengine 包中进行所有资源相关的watch机制  及时向cloudagent进行推送
package watchengine

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type WatcherEngine struct {
	Enable         bool
	MasterLBIp     string
	MasterLBPort   int
	RedisIp        string
	RedisPort      int
	PodInfo        *v1.Pod
	DeploymentInfo *appsv1.Deployment
}

func InitK8sClientWatchEngine(engin WatcherEngineInterface) {}
func init() {
	var watchEngin WatcherEngineInterface
	InitK8sClientWatchEngine(watchEngin)
}

func (we WatcherEngine) ListPods(namespace string) (unstructured.UnstructuredList, error) {
	fmt.Println(namespace)
	return unstructured.UnstructuredList{}, nil
}
