package tests

import (
	"fmt"
	devicemonitor "keep/edge/pkg/device_monitor"
	"testing"
)

func TestDeviceMonitor(t *testing.T) {
	fmt.Println("TestDeviceMonitor")
	monitor := devicemonitor.NewDeviceMonitor()
	defer monitor.Cleanup()

	monitor.Start()
	fmt.Println("TestDeviceMonitor END")
}
