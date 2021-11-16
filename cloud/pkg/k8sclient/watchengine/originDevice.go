package watchengine

import (
	"context"
	"fmt"
	"github.com/wonderivan/logger"
	"k8s.io/utils/pointer"
	native_engine "keep/cloud/pkg/k8sclient/native-engine"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreatePod(){

	clientset := native_engine.GetClient()
	//将配置信息赋值给deloymentClient
	deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	//构建deployment
	deployment:=&appsv1.Deployment{
		ObjectMeta:metav1.ObjectMeta{
			Name: "create-deplyment-test",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:                pointer.Int32Ptr(2),
			Selector:                &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":"demo",
				},
				MatchExpressions: nil,
			},
			Template:                apiv1.PodTemplateSpec{
				ObjectMeta:metav1.ObjectMeta{
					Labels: map[string]string{
						"app":"demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: "nginx",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name: "http",
									Protocol: apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},

						},
					},
				},
			},
		},
	}
	result, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err!=nil{
		panic(err)
	}
	logger.Info("create Pod Name:"+result.GetObjectMeta().GetName())
	fmt.Printf("Create Pod Name : %q \n",result.GetObjectMeta().GetName())
}
