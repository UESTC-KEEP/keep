package cloudimageManager

import (
	"keep/cloud/pkg/cloudimagemanager/config"
	"keep/cloud/pkg/common/modules"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/core"
	"keep/pkg/util/kplogger"
	"os"
)

type CloudImageManager struct {
	enable bool
}

func Register(cmi *cloudagent.CloudImageManager) {
	config.InitConfigure(cmi)
	rd, err := NewCloudImageManager(cmi.Enable)
	if err != nil {
		kplogger.Error("初始化RequestDispatcher失败...:", err)
		os.Exit(1)
	}
	core.Register(rd)
}

func (cim *CloudImageManager) Name() string {
	return modules.CloudImageManagerModule
}

func (cim *CloudImageManager) Group() string {
	return modules.CloudImageManagerGroup
}

func (cim *CloudImageManager) Start() {

}

func (cim *CloudImageManager) Cleanup() {}

func (cim *CloudImageManager) Enable() bool {
	return cim.enable
}

func NewCloudImageManager(enable bool) (*CloudImageManager, error) {
	return &CloudImageManager{
		enable: enable,
	}, nil
}
