package kafka

import (
	"github.com/Shopify/sarama"
	saramaCluster "github.com/bsm/sarama-cluster"
)

//var Address = []string{"192.168.1.103:9092", "192.168.1.103:9093"}

type ProducerMessage = sarama.ProducerMessage // bie ming
type ProducerError = sarama.ProducerError
type ConsumerMessage = sarama.ConsumerMessage
type NotifyMessage = saramaCluster.Notification
type ConsumerError = sarama.ConsumerError

type Config struct {
	saramaCluster.Config
	SyncProducerAmount    int
	AsyncProducerAmount   int
	ConsumerOfGroupAmount int
	OffsetLocalOrServer   int //0,local  1,server  2,newest
}

func NewConfig() (conf *Config) {
	conf = new(Config)
	conf.Config = *saramaCluster.NewConfig()
	conf.Config.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.Config.Producer.Return.Successes = true
	conf.Config.Producer.Return.Errors = true
	// conf.Config.Config.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.SyncProducerAmount = 1
	conf.AsyncProducerAmount = 1
	conf.ConsumerOfGroupAmount = 1
	conf.OffsetLocalOrServer = 1
	return
}
