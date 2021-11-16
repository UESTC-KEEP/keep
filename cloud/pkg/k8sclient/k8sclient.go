package k8sclient

<<<<<<< HEAD
import "keep/cloud/pkg/k8sclient/watchengine"

func main(){
	watchengine.CreatePod()
=======
import (
	"github.com/wonderivan/logger"
	v1 "k8s.io/api/core/v1"
	"keep/cloud/pkg/common"
	k8sclientconfig "keep/cloud/pkg/k8sclient/config"
	"keep/core"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"os"
	"time"
)

type K8sClient struct {
	enable       bool
	MasterLBIp   string
	MasterLBPort int
	PodInfo      *v1.Pod
}

// Register 注册healthzagent模块
func Register(k *cloudagent.K8sClient) {
	k8sclientconfig.InitConfigure(k)
	healthzagent, err := NewK8sClient(k.Enable)
	if err != nil {
		logger.Error("初始化k8sclient失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(healthzagent)
}

func (k *K8sClient) Name() string {
	return modules.K8sClientModule
}

func (k *K8sClient) Group() string {
	return modules.K8sClientGroup
>>>>>>> b0af266029c89d24fd39eac5960a66536ae9a802
}

//Enable indicates whether this module is enabled
func (k *K8sClient) Enable() bool {
	return k.enable
}

func (k *K8sClient) Start() {
	logger.Debug("k8sclient开始启动....")
	for {
		time.Sleep(time.Second)
	}
}

func NewK8sClient(enable bool) (*K8sClient, error) {
	return &K8sClient{
		enable: enable,
	}, nil
}
