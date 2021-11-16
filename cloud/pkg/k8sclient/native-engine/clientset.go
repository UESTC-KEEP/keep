package native_engine

import (
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"keep/cloud/pkg/k8sclient/conf"
)

func GetClient() *kubernetes.Clientset{
	config, err := conf.GetKubeConfig()
	if err!=nil{
		logger.Error(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err!=nil{
		logger.Error(err.Error())
	}
	return clientset
}
