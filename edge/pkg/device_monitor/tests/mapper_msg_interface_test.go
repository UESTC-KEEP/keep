package tests

import (
	"fmt"
	devicemonitor "keep/edge/pkg/device_monitor"
	"testing"
	"time"
)

//dummy mapper
func TestMapperMsgInterface(t *testing.T) {
	go TestDeviceMonitor(t)
	time.Sleep(time.Second)
	t.Log("TestMapperMsgInterface")
	fmt.Println("test")
	msg_interface := devicemonitor.NewMsgInterface("dummy")
	defer msg_interface.Destroy()

	for {
		time.Sleep(time.Second)
		msg := devicemonitor.NewMapperMsg()
		msg.AddItem("val", 31231231)
		msg.AddItem("name", "dhaskldak")
		msg_interface.SendStatusData(msg)
	}
}
