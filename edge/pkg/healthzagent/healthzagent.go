package healthzagent

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/wonderivan/logger"
	"keep/core"
	"keep/edge/pkg/common/modules"
	healthzagentconfig "keep/edge/pkg/healthzagent/config"
	"keep/edge/pkg/healthzagent/server"
	edgeagent "keep/pkg/apis/compoenentconfig/edgeagent/v1alpha1"
	"os"
	"time"
)

type HealthzAgent struct {
	enable             bool
	hostInfoStat       *host.InfoStat
	cpu                *[]cpu.InfoStat
	mem                *mem.VirtualMemoryStat
	diskPartitionStat  *[]disk.PartitionStat
	diskIOCountersStat *map[string]disk.IOCountersStat
	netIOCountersStat  *[]net.IOCountersStat
}

// Register 注册healthzagent模块
func Register(h *edgeagent.HealthzAgent) {
	healthzagentconfig.InitConfigure(h)
	healthzagent, err := NewHealthzAgent(h.Enable)
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
	logger.Debug(server.DescribeMachine(server.GetMachineStatus()))
	logger.Debug("healthzagent开始启动....")
	n := 0
	for {
		fmt.Println("healthzagent运行中：", n)
		n++
		time.Sleep(1 * time.Second)
	}
}

// NewHealthzAgent 创建新的healthzagent对象  并且初始化它
func NewHealthzAgent(enabled bool) (*HealthzAgent, error) {
	ha := &HealthzAgent{
		enable: enabled,
	}
	return ha, nil
}
