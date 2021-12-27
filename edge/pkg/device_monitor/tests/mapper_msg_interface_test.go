package tests

import (
	"fmt"
	devicemonitor "keep/edge/pkg/device_monitor"
	"testing"
)

func TestMapperMsgInterfacet(*testing.T) {
	fmt.Println("TestMapperMsgInterfacet")
	msg := devicemonitor.NewMsgInterface("dummy")
	defer msg.Destroy()
}
