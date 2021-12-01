package Router

import (
	"keep/cloud/pkg/common/kafka"
	"keep/pkg/util/core/model"
	logger "keep/pkg/util/loggerv1.0.1"
)

var p = kafka.NewProducerConfig("log")
var a = kafka.NewProducerConfig("add")

func MessageDispatcher(msg *model.Message) {

	switch msg.Router.Group {
	case "/log":
		sendToKafkaLog(msg)
	case "/add":
		sendToKafkaAdd(msg)
	}

}
func sendToKafkaLog(msg *model.Message) {
	kafkaMsg := msg.Content.(string)
	logger.Info("send to kafka", kafkaMsg)
	go func() { p.Msg <- kafkaMsg }()
}

func sendToKafkaAdd(msg *model.Message) {
	kafkaMsg := msg.Content.(string)
	logger.Info("send to kafka", kafkaMsg)
	go func() {
		a.Msg <- kafkaMsg
	}()
}

func ShutDownChan() {
	close(p.Msg)
	close(a.Msg)
}
