package edgetwin

import (
	"fmt"
	"keep/edge/pkg/common/modules"
	edgetwinconfig "keep/edge/pkg/edgetwin/config"
	"keep/edge/pkg/edgetwin/sqlite"
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"keep/pkg/util/core"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	"keep/pkg/util/loggerv1.0.1"
	"os"
	"time"
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
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("====================")
		testListPod()
	}()

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
	msg.Content = map[string]string{"namespace": "default"}
	beehiveContext.Send(modules.EdgePublisherModule, *msg)
}
