package tests

import (
	"context"
	"fmt"
	devicemonitor "keep/edge/pkg/device_monitor"
	"testing"
	"time"
)

func TestMapperMsgInterfacet(*testing.T) {
	fmt.Println("TestMapperMsgInterfacet")
	msg := devicemonitor.NewMsgInterface(context.Background(), "dummy")
	defer msg.Destroy()

	for {
		time.Sleep(time.Second)
		msg.SendStatusData([]byte("asdasdasdas\n"))
	}
}
