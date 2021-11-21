package logsagent

import (
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/logsagent/config"
	logsagentconfig "keep/edge/pkg/logsagent/config"
	"keep/edge/pkg/logsagent/tailf"
	"keep/pkg/util/core"

	"github.com/kubeedge/beehive/pkg/core"
	"github.com/wonderivan/logger"

	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"

	"os"
	"time"
)

type LogsAgent struct {
	enable      bool
	logLevel    int
	logTime     time.Time
	logFileName string
	logInfo     string
}

// Register 注册healthzagent模块
func Register(l *edgeagent.LogsAgent) {
	logsagentconfig.InitConfigure(l)
	healthzagent, err := NewLogsAgent(l.Enable)
	if err != nil {
		logger.Error("初始化logsagent失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(healthzagent)
}

func (l *LogsAgent) Cleanup() {
	//logger.Debug("准备清理模块：",modules.LogsAgentModule)
}

func (l *LogsAgent) Name() string {
	return modules.LogsAgentModule
}

func (l *LogsAgent) Group() string {
	return modules.LogsAgentGroup
}

//Enable indicates whether this module is enabled
func (l *LogsAgent) Enable() bool {
	return l.enable
}

func (l *LogsAgent) Start() {
	logger.Debug("logsagent开始启动....")
	logger.Info("所需监测日志文件:", config.Config.LogFiles)
	// 监听选中日志文件
	tailf.StartWatchingLogs(config.Config.LogFiles)
}

func NewLogsAgent(enable bool) (*LogsAgent, error) {
	return &LogsAgent{
		enable: enable,
	}, nil
}
