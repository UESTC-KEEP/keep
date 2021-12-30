package tests

import (
	"fmt"
	server "keep/edge/pkg/healthzagent/server"
	"testing"
)

func TestGetNetworkInterfaceInfoList(t *testing.T) {
	fmt.Println(server.GetNetworkInterfaceInfoList((&server.NetworkInterfaceInfoOption{
		Flags: true,
		MTU:   true})))
}
