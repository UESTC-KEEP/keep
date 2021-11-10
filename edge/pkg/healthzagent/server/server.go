/*
Package server 用于收集边缘主机的资源用量
*/
package server

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/wonderivan/logger"
	edgeagent "keep/pkg/apis/compoenentconfig/edgeagent/v1alpha1"
	"time"
)

// GetMachineStatus 获取节点综合信息
func GetMachineStatus() *edgeagent.HealthzAgent {
	var healagent edgeagent.HealthzAgent
	healagent.HostInfoStat, _ = GetBasicStatus()
	healagent.Cpu, _, _ = GetCpuStatus()
	healagent.Mem = GetMemStatus()
	healagent.DiskPartitionStat, healagent.DiskIOCountersStat = GetDiskStatus()
	healagent.NetIOCountersStat, _ = GetNetIOStatus()
	return &healagent
}

// DescribeMachine 描述主机信息
func DescribeMachine(ha *edgeagent.HealthzAgent) string {
	machine := (*ha).HostInfoStat
	hostname := machine.Hostname
	os := machine.OS
	KernelArch := machine.KernelArch
	KernelVersion := machine.KernelVersion
	PlatformFamily := machine.PlatformFamily
	_, cpupencent, _ := GetCpuStatus()
	return fmt.Sprintf("\n主机名：%s\n操作系统：%s\n内核架构：%s\n内核版本：%s\n平台族：%s\n内存用量：%.2f%%\ncpu用量：%.2f%%\n", hostname, os, KernelArch, KernelVersion, PlatformFamily, (*ha).Mem.UsedPercent, cpupencent)
}

// GetBasicStatus 获取边缘节点基础信息
func GetBasicStatus() (*host.InfoStat, error) {
	hInfo, err := host.Info()
	if err != nil {
		logger.Warn("获取边缘节点主机用量失败....", err)
		return nil, err
	}
	//logger.Debug("获取边缘节点基础信息:",hInfo)
	return hInfo, nil
}

// GetCpuStatus 获取当前边缘机器的cpu用量
func GetCpuStatus() (*[]cpu.InfoStat, float64, error) {
	cpuInfos, err := cpu.Info()
	if err != nil {
		logger.Warn("获取节点cpu信息失败：", err)
	}
	//logger.Debug("获取边缘节点的cpu信息：",cpuInfos)
	// CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
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

// GetDiskStatus 获取节点磁盘用量信息
func GetDiskStatus() (*[]disk.PartitionStat, *map[string]disk.IOCountersStat) {
	parts, err := disk.Partitions(true)
	if err != nil {
		logger.Warn("get Partitions failed, err:%v\n", err)
		return nil, nil
	}
	//fmt.Println(parts)
	for _, part := range parts {
		//fmt.Printf("part:%v\n", part.String())
		_, err := disk.Usage(part.Mountpoint)
		if err != nil {
			logger.Warn("路径：", part.Mountpoint, " 磁盘查询失败...", err)
			continue
		}
		//fmt.Printf("disk info:used:%v free:%v\n", diskInfo.UsedPercent, diskInfo.Free)
	}

	ioStat, _ := disk.IOCounters()
	//fmt.Println(ioStat)
	return &parts, &ioStat
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
