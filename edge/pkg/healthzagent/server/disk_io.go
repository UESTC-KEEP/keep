package server

import (
	"github.com/shirou/gopsutil/disk"
	logger "keep/pkg/util/loggerv1.0.1"
)

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
