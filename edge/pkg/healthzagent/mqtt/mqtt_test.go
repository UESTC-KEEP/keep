package mqtt

import (
	"fmt"
	"testing"
)

func TestSubscribeMqtt(t *testing.T) {
	SubscribeMqtt("192.168.1.40", "1883", "clock_sensor")
}

func TestGetRencentMqttMsg(t *testing.T) {
	fmt.Println(SubscribeMqtt("192.168.1.40", "1883", "clock_sensor"))
}
