package LogWatcher

import (
	"github.com/Shopify/sarama"
	"keep/cloud/pkg/common/kafka"
)

type LogStruct struct {
	Logid string `json:logid"`

}
var Producer sarama.SyncProducer

const(
	address="192.168.1.111"
	topic=""
)

func InitLogPusher(){
	config := sarama.NewConfig()
	Producer,_= sarama.NewSyncProducer([]string{address}, config)
}

func GetLogFromKafka(){
	messages := make(chan *sarama.ConsumerMessage)
	kafka.Subscribe([]string{address},topic,"",messages)
	PushertoKafka(messages)
}
func PushertoKafka(get chan *sarama.ConsumerMessage,){
	read:=<-get
	for i := 0; i < len(read.Key); i++ {
		msg:=&sarama.ProducerMessage{
			Topic: topic,
			Partition: int32(10),
			Key: sarama.ByteEncoder{read.Key[i]},
			Value: sarama.ByteEncoder{read.Value[i]},
		}
		Producer.SendMessage(msg)
	}

}


