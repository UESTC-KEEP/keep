package tests

import (
	devicemonitor "keep/edge/pkg/device_monitor"
	"testing"
	"time"
)

func TestMapperMsgInterfacet(t *testing.T) {
	t.Log("TestMapperMsgInterfacet")
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
