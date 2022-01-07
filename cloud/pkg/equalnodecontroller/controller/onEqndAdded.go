package controller

import (
	"context"
	crdv1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/equalnode/v1alpha1"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (eqndctl *EqualNodeController) equalNodeAdded(eqnd *crdv1.EqualNode) {
	eqndctl.equalnodeManager.EqualNode.Store(eqnd.Name, eqnd)
	logger.Info("----------- crd增加:  ")
	logger.Info("----------- crd增加:  ", eqnd)
	clientset := client.GetKubeClient()
	// 增加之后创建一个ngnix deplyment
	dep, err := clientset.AppsV1().Deployments("default").Create(context.Background(), &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "keep-eqnd-test-nginx",
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": eqnd.Spec.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": eqnd.Spec.Name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            "keep-eqnd-test-nginx",
							Image:           "nginx:1.13.5-alpine",
							ImagePullPolicy: "IfNotPresent",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1alpha1",
		},
	})
	if err != nil {
		logger.Error(err)
	}
	logger.Warn(dep)
	logger.Debug(eqnd.Spec.Name, "      ", eqnd.Spec.Eqnd)
}
