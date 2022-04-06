package server

import (
	"github.com/UESTC-KEEP/keep/pkg/util/kplogger"

	"github.com/shirou/gopsutil/disk"
)

// GetDiskStorageStatus 获取节点磁盘用量信息
func GetDiskStorageStatus() *map[string]*disk.UsageStat {
	parts, err := disk.Partitions(false) //只需要管实际的物理磁盘
	if err != nil {
		kplogger.Errorf("get Partitions failed, err:%v", err)
		return nil
	}

	usage_map := make(map[string]*disk.UsageStat)
	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint) //不是用Device获取分区空间信息
		if err != nil {
			kplogger.Errorf("Cannont access disk partition %s", part.Device)
			continue
		} else {
			usage_map[part.Device] = usage
		}
	}

	return &usage_map
}

func GetDiskIOStatus() *map[string]*disk.IOCountersStat {
	parts, err := disk.Partitions(false) //只需要管实际的物理磁盘
	if err != nil {
		kplogger.Warn("get Partitions failed, err:%v", err)
		return nil
	}

	io_map := make(map[string]*disk.IOCountersStat)
	for _, part := range parts {
		ioStat, err := disk.IOCounters(part.Device) //如果不指定分区名，那就会把根分区和子分区算在一起
		if err != nil {
			kplogger.Errorf("Cannont access disk partition %s", part.Device)
			continue
		} else {
			for _, v := range ioStat {
				io_map[part.Device] = &v
			}
		}
	}

	return &io_map
}
