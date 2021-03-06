package kafka

import (
	"github.com/UESTC-KEEP/keep/constants/cloud"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type ProducerConfig struct {
	Address []string
	Topic   string
	Msg     chan string
}

func NewProducerConfig(topic string) *ProducerConfig {
	msg := make(chan string)
	p := &ProducerConfig{
		Address: cloud.Address,
		// Address: []string{"192.168.1.103:9092", "192.168.1.103:9093"},
		Topic: topic,
		Msg:   msg,
	}
	return p
}

// Publish 向集群发送消息 producer
func (k *ProducerConfig) Publish() {
	config := NewConfig()

	produ, err := InitManualRetryAsyncProducer(k.Address, config)
	if err != nil {
		log.Println(err)
		return
	}
	defer produ.Close()

	go func(p *AsyncProducer) {
		for {
			select {
			case <-p.Successes():
				//log.Println("noerr")
			//fmt.Println("发送成功")
			//fmt.Println("offsetCfg:", suc.Offset, " partitions:", suc.Partition, " metadata:", suc.Metadata, " value:", value)
			case fail := <-p.Errors():
				log.Println("err: ", fail)
				// case <-p.Errors():
				// fmt.Println("no err")
			}
			time.Sleep(500 * time.Millisecond)
		}
	}(produ)

	for msgvalue := range k.Msg {
		//发送的消息,主题,key
		msg := &ProducerMessage{
			Topic: k.Topic,
		}
		//将字符串转化为字节数组
		msg.Value = sarama.ByteEncoder(msgvalue)
		//使用通道发送
		produ.Send() <- msg
		time.Sleep(500 * time.Millisecond)
	}
}
