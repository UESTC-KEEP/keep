package server

import (
	"context"
	"fmt"
	"os"

	net_io "github.com/shirou/gopsutil/net"
)

type NetworkInterfaceInfoOption struct { //控制输出网卡的哪些信息
	Index bool
	MTU   bool
	// Name         string 必须要有名称
	HardwareAddr bool
	Flags        bool
	Addrs        bool
}

type NetworkInterfaceInfo map[string]interface{}
type NetworkInterfaceInfoList []NetworkInterfaceInfo

func GetNetworkInterfaceInfoList(opt *NetworkInterfaceInfoOption) *NetworkInterfaceInfoList {
	net_interface, err := net_io.InterfacesWithContext(context.Background())
	if err != nil {
		fmt.Println("Failed to get network interface list")
		return nil
	}

	var list NetworkInterfaceInfoList

	for _, card := range net_interface {
		info := make(NetworkInterfaceInfo)
		info["Name"] = card.Name
		if opt.Index {
			info["Index"] = card.Index
		}
		if opt.Addrs {
			info["Addrs"] = card.Addrs
		}
		if opt.Flags {
			info["Flags"] = card.Flags
		}
		if opt.MTU {
			info["MTU"] = card.MTU
		}
		if opt.HardwareAddr {
			info["HardwareAddr"] = card.HardwareAddr
		}

		list = append(list, info)
	}

	return &list
}

// IsVirtualInterface 检测是否是虚拟网卡。当前仅限于linux系统
func IsVirtualInterface(interface_name string) (bool, error) {
	net_dev_path := "/sys/devices/virtual/net/" + interface_name
	_, err := os.Stat(net_dev_path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) { //如果不存在，则可能是物理设备，也可以是根本没这设备，所以interface_name需要是实际存在的设备名
		return false, nil
	} else {
		return false, err
	}
}

// type NetworkInterfaceIOcountOptions struct {
// 	BytesSent   bool
// 	BytesRecv   bool
// 	PacketsSent bool
// 	PacketsRecv bool
// 	Errin       bool
// 	Errout      bool
// 	Dropin      bool
// 	Dropout     bool
// 	Fifoin      bool
// 	Fifoout     bool
// }

// type NetworkInterfaceIOInfo map[string]uint64

// func NetworkInterfaceIOcount(interface_name string, opt *NetworkInterfaceIOcountOptions) *NetworkInterfaceIOInfo {
// 	info, err := net_io.IOCountersByFile(false, interface_name)
// 	if err != nil {
// 		fmt.Println("Failed to get IO info of ", interface_name)
// 		return nil
// 	}
// 	info_out := make(NetworkInterfaceIOInfo)

// 	if opt.BytesRecv {
// 		info_out["BytesRecv"] = info[0].BytesRecv
// 	}

// 	return &info_out
// }

// //这里不做筛选，因为要频繁调用，直接全部输出的效率比筛选要高
// func NetworkInterfaceIOcount(interface_name string) *net_io.IOCountersStat {
// 	info, err := net_io.IOCountersByFile(false, interface_name)
// 	if err != nil {
// 		kplogger.Error("Failed to get IO info of ", interface_name)
// 		return nil
// 	}
// 	return &(info[0])
// }
