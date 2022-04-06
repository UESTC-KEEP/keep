package edgetwin

import (
	"fmt"
	"github.com/UESTC-KEEP/keep/edge/pkg/common/modules"
	edgetwinconfig "github.com/UESTC-KEEP/keep/edge/pkg/edgetwin/config"
	"github.com/UESTC-KEEP/keep/edge/pkg/edgetwin/sqlite"
	edgeagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"github.com/UESTC-KEEP/keep/pkg/util/core"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/core/model"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"os"
)

type EdgeTwin struct {
	enable         bool
	sqliteFilePath string
}

// Register 注册healthzagent模块
func Register(et *edgeagent.EdgeTwin) {
	edgetwinconfig.InitConfigure(et)
	edgetwin, err := NewEdgeTwin(et.Enable)
	if err != nil {
		logger.Error("初始化EdgeTwin失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(edgetwin)
}

func (et *EdgeTwin) Cleanup() {

}

func (et *EdgeTwin) Name() string {
	return modules.EdgeTwinModule
}

func (et *EdgeTwin) Group() string {
	return modules.EdgeTwinGroup
}

//Enable indicates whether this module is enabled
func (et *EdgeTwin) Enable() bool {
	return et.enable
}

func (et *EdgeTwin) Start() {
	logger.Debug("EdgeTwin开始启动....")
	sqlite.ReceiveFromBeehiveAndInsert()
	// 测试边端向云端发送数据
	//go func() {
	//	time.Sleep(time.Second * 20)
	//	logger.Error("====================")
	//	testListPod()
	//}()

}

func NewEdgeTwin(enable bool) (*EdgeTwin, error) {
	return &EdgeTwin{
		enable: enable,
	}, nil
}

func testListPod() {
	msg := model.NewMessage("")
	msg.SetResourceOperation("$uestc/keep/k8sclient/naiveengine/pods/", "list")
	msg.SetRoute("$uestc/keep/k8sclient/naiveengine/pods/", "$uestc/keep/k8sclient/naiveengine/pods/")
	msg.Router.Source = modules.EdgeTwinModule
	reqMap := make(map[string]string)
	reqMap["namespace"] = "default"
	msg.Content = reqMap
	logger.Error(fmt.Sprintf("%#v", msg))
	beehiveContext.Send(modules.EdgePublisherModule, *msg)
}
