package healthzagent

import (
	"github.com/kubeedge/beehive/pkg/core"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/wonderivan/logger"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/healthzagent/config"
	prome "keep/edge/pkg/healthzagent/promethus"
	"keep/edge/pkg/healthzagent/server"

	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"

	"os"
)

type HealthzAgent struct {
	enable                    bool
	hostInfoStat              *host.InfoStat
	cpu                       *[]cpu.InfoStat
	mem                       *mem.VirtualMemoryStat
	diskPartitionStat         *[]disk.PartitionStat
	diskIOCountersStat        *map[string]disk.IOCountersStat
	netIOCountersStat         *[]net.IOCountersStat
	defaultEdgeHealthInterval int
	cpuUsage                  float64
}

// Register 注册healthzagent模块
func Register(h *edgeagent.HealthzAgent) {
	config.InitConfigure(h)
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
	logger.Debug("healthzagent开始启动....")
	// 打印机器配置
	server.GetMachineStatus()
	logger.Debug(server.DescribeMachine(&server.Healagent))
	// 启动周期性任务轮询本机用量
	cron := server.StartMetricEdgeInterval(config.Config.DefaultEdgeHealthInterval)
	// 启动本机StartMertricsServer
	go prome.StartMertricsServer(8080)
	//os.Exit(1)
	defer cron.Stop()
}

// NewHealthzAgent 创建新的healthzagent对象  并且初始化它
func NewHealthzAgent(enabled bool) (*HealthzAgent, error) {
	ha := &HealthzAgent{
		enable: enabled,
	}
	return ha, nil
}
