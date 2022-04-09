package docker

import (
	"context"
	"encoding/json"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"sync"
)

var clit *client.Client

// 初始化客户端
func init() {
	clit_, err := client.NewClientWithOpts()
	if err != nil {
		logger.Warn("初始化docker客户端失败...")
		return
	}
	clit = clit_
}

// GetAllDockerImages 查看节点镜像列表
func GetAllDockerImages() (*[]types.ImageSummary, error) {

	list, err := clit.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		logger.Warn(err)
		return nil, err
	}
	return &list, nil
}

func ListContainerMetrics() (*map[string]ContainerMetrics, error) {
	metricsMap := make(map[string]ContainerMetrics, 40)
	var mapLock sync.Mutex
	containers, err := clit.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, container := range containers {
		wg.Add(1)
		go func(cont types.Container) {
			defer wg.Done()
			resp, err := clit.ContainerStatsOneShot(context.Background(), cont.ID)
			if err != nil {
				logger.Error("get container stat error: ", err)
			}

			dec := json.NewDecoder(resp.Body)
			var (
				v                      *types.StatsJSON
				memPercent, cpuPercent float64
				blkRead, blkWrite      uint64 // Only used on Linux
				mem, memLimit          float64
				previousCPU            uint64
				previousSystem         uint64
				netRx, netTx           float64
			)

			if err = dec.Decode(&v); err != nil {
				logger.Error("decode container stat error: ", err)
			}
			previousCPU = v.PreCPUStats.CPUUsage.TotalUsage
			previousSystem = v.PreCPUStats.SystemUsage
			cpuPercent = calculateCPUPercentUnix(previousCPU, previousSystem, v)
			blkRead, blkWrite = calculateBlockIO(v.BlkioStats)
			mem = calculateMemUsageUnixNoCache(v.MemoryStats)
			memLimit = float64(v.MemoryStats.Limit)
			memPercent = calculateMemPercentUnixNoCache(memLimit, mem)
			netRx, netTx = calculateNetwork(v.Networks)

			mapLock.Lock()
			metricsMap[v.ID] = ContainerMetrics{v.Name, v.ID, cpuPercent, mem, memPercent, map[string]float64{"read": float64(blkRead), "write": float64(blkWrite)}, map[string]float64{"recv": netRx, "send": netTx}}
			mapLock.Unlock()
		}(container)
	}
	wg.Wait()
	return &metricsMap, nil
}

func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
		onlineCPUs  = float64(v.CPUStats.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(v.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}

func calculateBlockIO(blkio types.BlkioStats) (uint64, uint64) {
	var blkRead, blkWrite uint64
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		if len(bioEntry.Op) == 0 {
			continue
		}
		switch bioEntry.Op[0] {
		case 'r', 'R':
			blkRead = blkRead + bioEntry.Value
		case 'w', 'W':
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	return blkRead, blkWrite
}

func calculateMemUsageUnixNoCache(mem types.MemoryStats) float64 {
	// cgroup v1
	if v, isCgroup1 := mem.Stats["total_inactive_file"]; isCgroup1 && v < mem.Usage {
		return float64(mem.Usage - v)
	}
	// cgroup v2
	if v := mem.Stats["inactive_file"]; v < mem.Usage {
		return float64(mem.Usage - v)
	}
	return float64(mem.Usage)
}

func calculateMemPercentUnixNoCache(limit float64, usedNoCache float64) float64 {
	// MemoryStats.Limit will never be 0 unless the container is not running and we haven't
	// got any data from cgroup
	if limit != 0 {
		return usedNoCache / limit * 100.0
	}
	return 0
}

func calculateNetwork(network map[string]types.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}
