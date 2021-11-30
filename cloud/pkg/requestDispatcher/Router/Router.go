package Router

import (
	"keep/pkg/util/core/model"
	"keep/pkg/util/loggerv1.0.1"
	"sync"
)

var Address = []string{"192.168.1.103:9092", "192.168.1.103:9093"}
var msgChan = make(chan string, 100)
var kafkaOnce sync.Once

func MessageDispatcher(msg *model.Message) {
	kafkaOnce.Do(func() {
		//kafka.AsyncProducer(Address, "topic", msgChan)
	})

	kafkaMsg := msg.Content.(string)
	logger.Info("send to kafka", kafkaMsg)
	msgChan <- kafkaMsg

}
