/*
Package server 用于收集边缘主机的资源用量
*/
package server

import (
	"fmt"
	"github.com/UESTC-KEEP/keep/edge/pkg/common/modules"
	"github.com/UESTC-KEEP/keep/edge/pkg/healthzagent/config"
	edgeagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/core/model"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"github.com/robfig/cron"
	"github.com/shirou/gopsutil/host"
	"strconv"
)

var Healagent edgeagent.HealthzAgent

type EdgeMetrics struct {
	MemoryIntototal  int `json:"memoryIntototal"`
	MemInuseMiBytes  int `json:"memInuseMiBytes"`
	MemInuserPercent int `json:"memInuserPercent"`

	CpuIntotalMi    int `json:"cpuIntotal"`
	CpuInusePercent int `json:"cpuInusePercent"`

	DiskIntotalMiBytes   int              `json:"diskIntotal"`
	DiskInuseByPartition []map[string]int `json:"diskInuseByPartition"`
	DiskIO               int              `json:"diskIO"`

	NetBandWidth int `json:"netBandWidth"`
	NetIO        int `json:"netIO"`
}

// GetMachineStatus  获取节点综合信息
func GetMachineStatus() {
	logger.Trace("查询节点用量...")
	Healagent.HostInfoStat, _ = GetBasicStatus()
	Healagent.Cpu, Healagent.CpuUsage, _ = GetCpuStatus()
	Healagent.Mem = GetMemStatus()
	Healagent.DiskPartitionStat = GetDiskStorageStatus()
	Healagent.DiskIOCountersStat = GetDiskIOStatus()
	Healagent.NetIOCountersStat, _ = GetNetIOStatus()
	Healagent.NodeName = config.Config.NodeName

	//resp, err := http.Get("127.0.0.1:9200/metrics")
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//同步数据到sqlite
	msg := model.Message{
		Header: model.MessageHeader{},
		Router: model.MessageRoute{
			Source:    modules.HealthzAgentModule,
			Group:     "",
			Operation: "",
			Resource:  "uestc-keep-kafka-metrics",
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
	logger.Trace(interval, "s 定时检测边缘节点健康状态启动成功...")
	return cron2
}
