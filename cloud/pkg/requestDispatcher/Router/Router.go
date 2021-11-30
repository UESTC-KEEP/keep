package Router

import (
	"github.com/wonderivan/logger"
	"keep/cloud/pkg/common/kafka"
	"keep/pkg/util/core/model"
	"sync"
)

var Address = []string{"192.168.1.103:9092", "192.168.1.103:9093"}
var msgChan = make(chan string, 100)
var kafkaOnce sync.Once

func MessageDispatcher(msg *model.Message) {
	kafkaOnce.Do(func() {
		kafka.AsyncPro(Address, "topic", msgChan)
	})

	kafkaMsg := msg.Content.(string)
	logger.Info("send to kafka", kafkaMsg)
	msgChan <- kafkaMsg

}
