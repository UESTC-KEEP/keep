// Package watchengine 包中进行所有资源相关的watch机制  及时向cloudagent进行推送
package watchengine

import (
	"github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/kubeedge-engine/devices/informer"
	naive_engine "github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/naive-engine"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
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

func StartAllInformers() {
	go naive_engine.StartNaiveEngineInformers()
	go device_informer.StartDeviceInformer()
}

func InitK8sClientWatchEngine(engin WatcherEngineInterface) {

}

func init() {
	var watchEngin WatcherEngineInterface
	InitK8sClientWatchEngine(watchEngin)
	//initK8sClientWatchEngine(watchEngin)
}
