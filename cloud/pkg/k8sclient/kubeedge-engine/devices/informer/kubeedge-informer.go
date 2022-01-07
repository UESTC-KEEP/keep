package device_informer

import (
	"fmt"
	kedevice "github.com/UESTC-KEEP/keep/cloud/pkg/apis/kubeedge/devices/v1alpha2"
	"github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/config"
	kubeedge_engine "github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/kubeedge-engine"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"time"
)

func StartDeviceInformer() {
	logger.Debug("启动device ingfoinformer...")
	informer := dynamicinformer.NewFilteredDynamicInformer(config.DynamicClient, kubeedge_engine.KubeedgeGVR, "default", 2*time.Second, nil, nil).Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    OnAddDevice,
		UpdateFunc: nil,
		DeleteFunc: nil,
	})
	informer.Run(beehiveContext.Done())
}

func OnAddDevice(newObj interface{}) {
	// 类型转换
	device := &kedevice.Device{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(newObj.(*unstructured.Unstructured).UnstructuredContent(), device)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug("informer:增加device....")
	fmt.Println(device)
}
