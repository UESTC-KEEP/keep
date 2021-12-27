package tests

import (
	"fmt"
	devicemonitor "keep/edge/pkg/device_monitor"
	"testing"
)

func TestDeiceMonitor(t *testing.T) {
	fmt.Println("TestDeiceMonitor")
	monitor := devicemonitor.NewDeviceMonitor()
	defer monitor.Destroy()

	monitor.Run()
}
