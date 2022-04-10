package kafka

import (
	"github.com/Shopify/sarama"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
)

//var Address = []string{"192.168.1.103:9092", "192.168.1.103:9093"}

type AsyncProducer struct {
	producer        sarama.AsyncProducer
	Id              int
	ProducerGroupId string
}

func InitManualRetryAsyncProducer(addr []string, conf *Config) (*AsyncProducer, error) {
	// conf.Config.Producer.Retry.Max = 0  默认为3

	aSyncProducer := &AsyncProducer{
		Id:              0,
		ProducerGroupId: "",
	}
	var err error
	aSyncProducer.producer, err = sarama.NewAsyncProducer(addr, &conf.Config.Config)
	if err != nil {
		logger.Info("sarama.NewSyncProducer err, message="+err.Error()+" ip:", addr)
		return nil, err
	} else {
		logger.Info("kafka集群链接成功：", addr)
	}
	return aSyncProducer, nil
}

//send message
func (asp *AsyncProducer) Send() chan<- *ProducerMessage {
	return asp.producer.Input()
}

func (asp *AsyncProducer) Successes() <-chan *ProducerMessage {
	return asp.producer.Successes()
}

func (asp *AsyncProducer) Errors() <-chan *ProducerError {
	return asp.producer.Errors()
}

func (asp *AsyncProducer) Close() (err error) {
	err = asp.producer.Close()
	return
}
