package lister

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kedevice "keep/cloud/pkg/apis/kubeedge/devices/v1alpha2"
	"keep/cloud/pkg/k8sclient/config"
	kubeedge_engine "keep/cloud/pkg/k8sclient/kubeedge-engine"
	logger "keep/pkg/util/loggerv1.0.1"
)

func GetAllDevice() (*kedevice.DeviceList, error) {
	devicelist, err := config.DynamicClient.Resource(kubeedge_engine.KubeedgeGVR).Namespace(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	// 实例化devicelist
	deviceList := &kedevice.DeviceList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(devicelist.UnstructuredContent(), deviceList)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return deviceList, err
}
