/*
Package server 用于收集边缘主机的资源用量
*/
package server

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"keep/edge/pkg/common/modules"
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	"keep/pkg/util/loggerv1.0.1"
	"strconv"
	"time"
)

var Healagent edgeagent.HealthzAgent

// GetMachineStatus  获取节点综合信息
func GetMachineStatus() {
	logger.Debug("查询节点用量...")
	Healagent.HostInfoStat, _ = GetBasicStatus()
	Healagent.Cpu, Healagent.CpuUsage, _ = GetCpuStatus()
	Healagent.Mem = GetMemStatus()
	Healagent.DiskPartitionStat, Healagent.DiskIOCountersStat = GetDiskStatus()
	Healagent.NetIOCountersStat, _ = GetNetIOStatus()
	//同步数据到sqlite
	msg := model.Message{
		Header: model.MessageHeader{},
		Router: model.MessageRoute{
			Source:    modules.HealthzAgentModule,
			Group:     "",
			Operation: "",
			Resource:  "",
		},
		Content: Healagent,
	}
	beehiveContext.Send(modules.EdgePublisherModule, msg)

	//logger.Debug(fmt.Sprintf("\n内存用量：%.2f%%  cpu用量：%.2f%% ", Healagent.Mem.UsedPercent, Healagent.CpuUsage))
	//for i := 0; i < len(*Healagent.NetIOCountersStat); i++ {
	//	fmt.Println(fmt.Sprintf("网卡名字：%s 发送数据：%vKB 接收数据：%vKB",
	//		(*Healagent.NetIOCountersStat)[i].Name, (*Healagent.NetIOCountersStat)[i].BytesSent/1024, (*Healagent.NetIOCountersStat)[i].BytesRecv/1024))
	//}
}

// DescribeMachine 描述主机信息
func DescribeMachine(ha *edgeagent.HealthzAgent) string {
	machine := (*ha).HostInfoStat
	hostname := machine.Hostname
	os := machine.OS
	KernelArch := machine.KernelArch
	KernelVersion := machine.KernelVersion
	PlatformFamily := machine.PlatformFamily
	return fmt.Sprintf("\n主机名：%s\n操作系统：%s\n内核架构：%s\n内核版本：%s\n平台族：%s\n内存用量：%.2f%%\ncpu用量：%.2f%%\n", hostname, os, KernelArch, KernelVersion, PlatformFamily, (*ha).Mem.UsedPercent, ha.CpuUsage)
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
			//logger.Warn("路径：", part.Mountpoint, " 磁盘查询失败...", err)
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

// StartMetricEdgeInterval  定时轮询边缘节点用量
func StartMetricEdgeInterval(interval int) *cron.Cron {
	//定期 删除过期文件
	cron2 := cron.New() //创建一个cron实例
	//validTime秒执行一次 更新边缘节点信息
	err := cron2.AddFunc("*/"+strconv.Itoa(interval)+" * * * * *", GetMachineStatus)
	if err != nil {
		logger.Error(err)
	}
	cron2.Start()
	logger.Debug("定时检测边缘节点健康状态启动成功...")
	return cron2
}
