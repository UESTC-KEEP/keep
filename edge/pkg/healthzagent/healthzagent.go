package healthzagent

import (
	"fmt"
	"github.com/wonderivan/logger"
	"keep/core"
	"keep/edge/pkg/common/modules"
	healthzagentconfig "keep/edge/pkg/healthzagent/config"
	edgeagent "keep/pkg/apis/compoenentconfig/edgeagent/v1alpha1"
	"os"
	"time"
)

type HealthzAgent struct {
	enable bool
	cpu    float64
	disk   float64
}

// Register 注册healthzagent模块
func Register(h *edgeagent.HealthzAgent) {
	healthzagentconfig.InitConfigure(h)
	healthzagent, err := newHealthzAgent(h.Enable)
	if err != nil {
		logger.Error("初始化edgehealthzagent失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(healthzagent)
}

func (h *HealthzAgent) Name() string {
	return modules.HealthzAgentModule
}

func (h *HealthzAgent) Group() string {
	return modules.HealthzAgentGroup
}

//Enable indicates whether this module is enabled
func (h *HealthzAgent) Enable() bool {
	return h.enable
}

func (h *HealthzAgent) Start() {
	logger.Debug("healthzagent开始启动....")
	n := 0
	for {
		fmt.Println("healthzagent运行中：", n)
		n++
		time.Sleep(1 * time.Second)
	}
}

// 创建新的healthzagent对象  并且初始化它
func newHealthzAgent(enabled bool) (*HealthzAgent, error) {
	ha := &HealthzAgent{
		enable: enabled,
		cpu:    healthzagentconfig.Config.Cpu,
		disk:   healthzagentconfig.Config.Disk,
	}
	return ha, nil
}
