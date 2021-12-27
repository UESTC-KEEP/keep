package tests

import (
	"fmt"
	devicemonitor "keep/edge/pkg/device_monitor"
	"testing"
	"time"
)

func TestMapperMsgInterfacet(*testing.T) {
	fmt.Println("TestMapperMsgInterfacet")
	msg := devicemonitor.NewMsgInterface("dummy")
	defer msg.Destroy()

	for {
		time.Sleep(time.Second)
		msg := NewMapperMsg()
		msg.SendStatusData(msg)
	}
}
