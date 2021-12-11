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
	"keep/constants/edge"
	"keep/pkg/util"
	logger "keep/pkg/util/loggerv1.0.1"
	"os"
	"strings"
	"time"
)

// NewDefaultEdgeAgentConfig returns a full EdgeCoreConfig object
func NewDefaultEdgeAgentConfig() *EdgeAgentConfig {
	hostnameOverride, _ := os.Hostname()
	localIP, _ := util.GetLocalIP(hostnameOverride)
	logger.Info(localIP)
	return &EdgeAgentConfig{
		DataBase: &DataBase{
			DataSource: DataBaseDataSource,
		},
		Modules: &Modules{
			HealthzAgent: &HealthzAgent{
				Enable:                    true,
				CpuUsage:                  0.0,
				DefaultEdgeHealthInterval: edge.DefaultEdgeHealthInterval,
				Cpu:                       nil,
				Mem:                       nil,
				DiskPartitionStat:         nil,
				DiskIOCountersStat:        nil,
				NetIOCountersStat:         nil,
				DeviceMqttTopics:          strings.Split(edge.DefaultDeviceMqttTopics, ";"),
			},
			LogsAgent: &LogsAgent{
				Enable:      true,
				LogLevel:    6,
				LogTime:     time.Now(),
				LogFileName: "",
				LogInfo:     "",
				LogFiles:    []string{edge.DefaultEdgeLogFiles},
			},
			EdgePublisher: &EdgePublisher{
				Enable:            true,
				HTTPServer:        edge.DefaultHttpServer,
				Port:              edge.DefaultCloudHttpPort,
				ServePort:         edge.DefaultEdgePort,
				Heartbeat:         edge.DefaultEdgeHeartBeat,
				EdgeMsgQueens:     []string{edge.DefaultLogsTopic, edge.DefaultDataTopic},
				TLSCAFile:         "",
				TLSCertFile:       "",
				TLSPrivateKeyFile: "",
				BeehiveTimeout:    edge.DefaultBeehiveTimeout * time.Millisecond,
				HostnameOverride:  hostnameOverride,
				LocalIP:           localIP,
				Token:             "",
			},
			EdgeTwin: &EdgeTwin{
				Enable:         true,
				SqliteFilePath: edge.DefaultEdgeTwinSqliteFilePath,
				BeehiveTimeout: edge.DefaultBeehiveTimeout * time.Millisecond,
			},
			DeviceMapperInterface: &DeviceMapperInterface{
				Enable: true,
			},
		},
	}
}

// NewMinEdgeCoreConfig returns a common EdgeCoreConfig object
func NewMinEdgeCoreConfig() *EdgeAgentConfig {
	hostnameOverride, _ := os.Hostname()
	localIP, _ := util.GetLocalIP(hostnameOverride)
	logger.Info(localIP)
	return &EdgeAgentConfig{
		DataBase: &DataBase{
			DataSource: DataBaseDataSource,
		},
		Modules: &Modules{
			HealthzAgent: &HealthzAgent{
				Enable:             true,
				Cpu:                nil,
				Mem:                nil,
				DiskPartitionStat:  nil,
				DiskIOCountersStat: nil,
				NetIOCountersStat:  nil,
			},
		},
	}
}
