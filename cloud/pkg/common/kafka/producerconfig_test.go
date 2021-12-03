package kafka

import (
	"testing"
	"time"
)

func TestProducerConfig_Publish(t *testing.T) {
	p := NewProducerConfig("topic")
	go p.Publish()

	// 测试push数据
	for i := 0;;i++ {
		time11 := time.Now()
		value := "this is a message 0805 "+time11.Format("15:04:05")
		p.Msg <- value
	}
}