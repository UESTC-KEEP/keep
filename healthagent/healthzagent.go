package main

import (
	"keep/constants/edge"
	"keep/healthagent/config"
	prome "keep/healthagent/promethus"
	"keep/healthagent/server"
	"keep/pkg/util/kplogger"
	logger "keep/pkg/util/loggerv1.0.1"
	"strings"

	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"os"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
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

func (h *HealthzAgent) Cleanup() {
	//logger.Debug("准备清理模块：",modules.HealthzAgentModule)
	prome.StopMertricsServer()
}

func (h *HealthzAgent) Start() {
	logger.Debug("healthzagent开始启动....")
	// 打印机器配置
	server.GetMachineStatus()
	kplogger.Debug(server.DescribeMachine(&server.Healagent))
	// 启动周期性任务轮询本机用量
	//cron := server.StartMetricEdgeInterval(config.Config.DefaultEdgeHealthInterval)
	// 启动本机StartMertricsServer
	server.StartMetricEdgeInterval(3)
	go prome.StartMertricsServer(edge.DefaultMetricsPort)
	//os.Exit(1)
	//defer cron.Stop()
}

// NewHealthzAgent 创建新的healthzagent对象  并且初始化它
func NewHealthzAgent(enabled bool) (*HealthzAgent, error) {
	ha := &HealthzAgent{
		enable: enabled,
	}
	return ha, nil
}

func main() {

	h := edgeagent.HealthzAgent{
		Enable:                    true,
		CpuUsage:                  0.0,
		DefaultEdgeHealthInterval: edge.DefaultEdgeHealthInterval,
		Cpu:                       nil,
		Mem:                       nil,
		DiskPartitionStat:         nil,
		DiskIOCountersStat:        nil,
		NetIOCountersStat:         nil,
		DeviceMqttTopics:          strings.Split(edge.DefaultDeviceMqttTopics, ";"),
	}

	config.InitConfigure(&h)

	healthzagent, err := NewHealthzAgent(h.Enable)
	if err != nil {
		logger.Error("初始化edgehealthzagent失败...:", err)
		os.Exit(1)
		return
	}

	defer healthzagent.Cleanup()
	healthzagent.Start()
}
