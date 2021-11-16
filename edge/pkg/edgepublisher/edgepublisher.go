package edgepublisher

import (
	"github.com/wonderivan/logger"
	"keep/core"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgepublisher/bufferpooler"
	"keep/edge/pkg/edgepublisher/chanmsgqueen"
	edgepublisherconfig "keep/edge/pkg/edgepublisher/config"
	"keep/edge/pkg/edgepublisher/publisher"

	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"

	"os"
	"sync"
)

type EdgePublisher struct {
	enable            bool
	httpserver        string
	port              int32
	heartbeat         int32
	tlscafile         string
	tlscertfile       string
	tlsprivatekeyfile string
	edgemsgqueens     []string
}

// Register 注册healthzagent模块
func Register(ep *edgeagent.EdgePublisher) {
	edgepublisherconfig.InitConfigure(ep)
	edgepublisher, err := NewEdgePublisher(ep.Enable)
	if err != nil {
		logger.Error("初始化logsagent失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(edgepublisher)
}

func (ep *EdgePublisher) Name() string {
	return modules.EdgePublisherModule
}

func (ep *EdgePublisher) Group() string {
	return modules.EdgePublisherGroup
}

//Enable indicates whether this module is enabled
func (ep *EdgePublisher) Enable() bool {
	return ep.enable
}

func (l *EdgePublisher) Start() {
	var wg sync.WaitGroup
	logger.Debug("EdgePublisher 开始启动....")
	// 启动边端服务20350
	// 初始化队列 确保队列初始化完成再启动服务
	chanmsgqueen.InitMsgQueens()
	wg.Wait()
	go bufferpooler.StartEdgePublisher()
	publisher.ReadQueueAndPublish()
}

func NewEdgePublisher(enable bool) (*EdgePublisher, error) {
	return &EdgePublisher{
		enable: enable,
	}, nil
}
