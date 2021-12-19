/*
Copyright 2021 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	crdClientset "keep/cloud/pkg/client/clientset/versioned"
	"os"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	cloudagentConfig "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
)

var keClient *kubeEdgeClient
var once sync.Once

func InitKubeEdgeClient(config *cloudagentConfig.K8sClient) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("",
		config.KubeAPIConfig.KubeConfig)
	if err != nil {
		klog.Errorf("Failed to build config, err: %v", err)
		os.Exit(1)
	}
	kubeConfig.QPS = float32(config.KubeAPIConfig.QPS)
	kubeConfig.Burst = int(config.KubeAPIConfig.Burst)

	dynamicClient := dynamic.NewForConfigOrDie(kubeConfig)

	kubeConfig.ContentType = runtime.ContentTypeProtobuf
	kubeClient := kubernetes.NewForConfigOrDie(kubeConfig)

	crdKubeConfig := rest.CopyConfig(kubeConfig)
	crdKubeConfig.ContentType = runtime.ContentTypeJSON
	crdClient := crdClientset.NewForConfigOrDie(crdKubeConfig)

	once.Do(func() {
		keClient = &kubeEdgeClient{
			kubeClient:    kubeClient,
			crdClient:     crdClient,
			dynamicClient: dynamicClient,
		}
	})
}

func GetKubeClient() kubernetes.Interface {
	return keClient.kubeClient
}

func GetCRDClient() crdClientset.Interface {
	return keClient.crdClient
}

func GetDynamicClient() dynamic.Interface {
	return keClient.dynamicClient
}

type kubeEdgeClient struct {
	kubeClient    *kubernetes.Clientset
	crdClient     *crdClientset.Clientset
	dynamicClient dynamic.Interface
}
