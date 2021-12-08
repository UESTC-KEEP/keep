package kafka

import (
	"fmt"
	"testing"
)

func TestConsumerConfig_Subscribe(t *testing.T) {
	c := NewConsumerConfig("topic1", "groupid")
	go c.Subscribe()
	for msg := range c.Ans{
		fmt.Println(string(msg.Value))
	}
}
