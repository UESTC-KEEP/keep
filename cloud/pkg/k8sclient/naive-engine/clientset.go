package naive_engine

import (
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"keep/constants"
)

func GetClient() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", constants.DefaultKubeConfigPath)
	if err != nil {
		logger.Error(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err.Error())
	}
	return clientset
}
