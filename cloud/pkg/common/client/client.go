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
	eqndcrdClientset "github.com/UESTC-KEEP/keep/cloud/pkg/client/eqnd/clientset/versioned"
	tenantClientset "github.com/UESTC-KEEP/keep/cloud/pkg/client/tenant/clientset/versioned"
	trqcrdClientset "github.com/UESTC-KEEP/keep/cloud/pkg/client/trq/clientset/versioned"
	"os"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	cloudagentConfig "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
)

var kpClient *kubeEdgeClient
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
	eqndcrdClient := eqndcrdClientset.NewForConfigOrDie(crdKubeConfig)
	trqcrdClient := trqcrdClientset.NewForConfigOrDie(crdKubeConfig)
	tenantClient := tenantClientset.NewForConfigOrDie(crdKubeConfig)

	once.Do(func() {
		kpClient = &kubeEdgeClient{
			kubeClient:    kubeClient,
			eqndcrdClient: eqndcrdClient,
			trqcrdClient:  trqcrdClient,
			tenantClient:  tenantClient,
			dynamicClient: dynamicClient,
		}
	})
}

func GetKubeClient() kubernetes.Interface {
	return kpClient.kubeClient
}

func GetEqndCRDClient() eqndcrdClientset.Interface {
	return kpClient.eqndcrdClient
}

func GetTrqCRDClient() trqcrdClientset.Interface {
	return kpClient.trqcrdClient
}

func GetTenantClient() tenantClientset.Interface {
	return kpClient.tenantClient
}

func GetDynamicClient() dynamic.Interface {
	return kpClient.dynamicClient
}

type kubeEdgeClient struct {
	kubeClient    *kubernetes.Clientset
	eqndcrdClient *eqndcrdClientset.Clientset
	trqcrdClient  *trqcrdClientset.Clientset
	tenantClient  *tenantClientset.Clientset
	dynamicClient dynamic.Interface
}
