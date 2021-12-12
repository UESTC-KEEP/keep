package Deployments

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"keep/cloud/pkg/k8sclient/config"
	logger "keep/pkg/util/loggerv1.0.1"
)

func GetDeploymentByName(name string) (*v1.Deployment, error) {
	deployment, err := config.Clientset.AppsV1().Deployments(metav1.NamespaceAll).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return deployment, err
}
