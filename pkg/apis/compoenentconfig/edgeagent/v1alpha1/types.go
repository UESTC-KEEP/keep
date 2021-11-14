/*
Copyright 2019 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"
)

const (
	// DataBaseDriverName is sqlite3
	DataBaseDriverName = "sqlite3"
	// DataBaseAliasName is default
	DataBaseAliasName = "default"
	// DataBaseDataSource is edge.db
	DataBaseDataSource = "/var/lib/keepedge/edgeagent.db"
)

// EdgeAgentConfig 是从edgeagent的配置文件读出来的内容
type EdgeAgentConfig struct {
	// 数据库信息
	DataBase *DataBase `json:"database,omitempty"`
	// edgeagentagent中的模块配置
	Modules *Modules `json:"modules,omitempty"`
}

// DataBase 说明数据库信息
type DataBase struct {
	// DriverName 数据库驱动
	// default "sqlite3"
	DriverName string `json:"driverName,omitempty"`
	// AliasName 别名
	// default "default"
	AliasName string `json:"aliasName,omitempty"`
	// DataSource 数据原
	// default "/var/lib/keepedge/edgeagent.db"
	DataSource string `json:"dataSource,omitempty"`
}

type Modules struct {
	HealthzAgent *HealthzAgent `json:"healthzAgent,omitempty"`
	LogsAgent    *LogsAgent    `json:"logsagent,omitempty"`
}

// HealthzAgent 是该模块的说明
// healthzagent用于收集当前边缘节点的用量信息  包括 网络带宽  磁盘用量 cpu 内存用量
type HealthzAgent struct {
	// Enable 说明healthzagent模块当前是否启用 如果没有启用则其对应的配置文件也不会进行校验 默认启动
	// default true
	Enable                    bool                            `json:"enable,omitempty"`
	HostInfoStat              *host.InfoStat                  `json:"host_info_stat"`
	Cpu                       *[]cpu.InfoStat                 `json:"cpu,omitempty"`
	CpuUsage                  float64                         `json:"cpu_usage,omitempty"`
	Mem                       *mem.VirtualMemoryStat          `json:"mem,omitempty"`
	DiskPartitionStat         *[]disk.PartitionStat           `json:"disk_partition_stat,omitempty"`
	DiskIOCountersStat        *map[string]disk.IOCountersStat `json:"disk_io_counters_stat"`
	NetIOCountersStat         *[]net.IOCountersStat           `json:"net_io_counters_stat"`
	DefaultEdgeHealthInterval int                             `json:"defaultEdgeHealthInterval,omitempty"`
	DefaultMqttCacheQueueSize int                             `json:"defaultMqttCacheQueueSize,omitempty"`
}

type LogsAgent struct {
	Enable      bool      `json:"enable,omitempty"`
	LogLevel    int       `json:"log_level"`
	LogTime     time.Time `json:"log_time"`
	LogFileName string    `json:"log_file_name"`
	LogInfo     string    `json:"log_info"`
	LogFiles    []string  `json:"log_files"`
}
