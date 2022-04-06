package server

import (
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// GetCpuStatus 获取当前边缘机器的cpu用量
func GetCpuStatus() (*[]cpu.InfoStat, float64, error) {
	cpuInfos, err := cpu.Info()
	if err != nil {
		logger.Warn("获取节点cpu信息失败：", err)
	}
	//logger.Debug("获取边缘节点的cpu信息：",cpuInfos)
	// CPU使用率
	percent, err := cpu.Percent(time.Second, false)
	return &cpuInfos, percent[0], err
}

// GetMemStatus 获取节点的内存使用情况
func GetMemStatus() *mem.VirtualMemoryStat {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		logger.Warn("获取节点内存信息失败：", err)
	}
	//logger.Debug("获取边缘节点的mem信息：",memInfo)
	return memInfo
}

// GetNetIOStatus 获取边缘节点网络使用情况
func GetNetIOStatus() (*[]net.IOCountersStat, error) {
	info, err := net.IOCounters(true)
	if err != nil {
		logger.Warn("网络流量查询失败...", err)
		return nil, err
	}
	//logger.Debug("网络信息：",info)
	return &info, nil
}
