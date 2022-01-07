package kubeedge_engine

import (
	kedevice "github.com/UESTC-KEEP/keep/cloud/pkg/apis/kubeedge/devices/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var KubeedgeGVR = schema.GroupVersionResource{
	Group:    "devices.kubeedge.io",
	Version:  "v1alpha2",
	Resource: "devices",
}

type KubeedgeEngine interface {
	// GetDevicesByNodeName 查询某个节点上的设备信息
	/*
		nodename: 节点名称
	*/
	GetDevicesByNodeName(nodename string) (*kedevice.DeviceList, error)
}
