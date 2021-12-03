package controller

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformer "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
	"time"
)

func TestController(t *testing.T) {
	clientconfig, _ := clientcmd.BuildConfigFromFlags("", "/home/et/.kube/config")
	clientset := kubernetes.NewForConfigOrDie(clientconfig)
	var nodes *v1.NodeList

	nodes, _ = clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		Watch: false,
	})
	//clientset.CoreV1().
	fmt.Println(nodes)
	//go func(){
	//
	//	time.Sleep(100*time.Millisecond)
	//}()
}

func TestInformer(t *testing.T) {
	clientconfig, _ := clientcmd.BuildConfigFromFlags("", "/home/et/.kube/config")
	clientset := kubernetes.NewForConfigOrDie(clientconfig)
	kubeinformerFactory := kubeinformer.NewSharedInformerFactory(clientset, 2000000000)
	deployinformer := kubeinformerFactory.Apps().V1().Deployments()
	deployinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    AddDeployment,
		UpdateFunc: nil,
		DeleteFunc: nil,
	})

	replicasinformer := kubeinformerFactory.Apps().V1().ReplicaSets().Informer()
	replicasinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    AddReplica,
		UpdateFunc: nil,
		DeleteFunc: nil,
	})

	podinformer := kubeinformerFactory.Core().V1().Pods().Informer()
	podinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    AddPod,
		UpdateFunc: nil,
		DeleteFunc: nil,
	})

	stopchan, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))

	go kubeinformerFactory.Start(stopchan.Done())
	defer cancel()
	for {
		select {
		case <-stopchan.Done():
			cancel()
		default:
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func AddDeployment(newObj interface{}) {
	fmt.Println("=======================================================================================")
	fmt.Println(newObj)
}

func AddReplica(newObj interface{}) {
	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println(newObj)
}

func AddPod(newObj interface{}) {
	fmt.Println("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	fmt.Println(newObj)
}
