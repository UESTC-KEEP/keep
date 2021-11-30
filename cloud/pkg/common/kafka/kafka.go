package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

type ConsumerConfig struct {
	Address []string
	Topic string
	GroupId string
	Ans chan *ConsumerMessage
}

type ProducerConfig struct {
	Address []string
	Topic string
	Msg  chan string
}

// Subscribe 订阅/接收消息 consumer
func (k *ConsumerConfig) Subscribe() error {
	config :=NewConfig()
	// ans := make(chan *sarama.ConsumerMessage)
	c , err := InitOneConsumerOfGroup(k.Address, k.Topic, k.GroupId, config)
	if err!=nil{
		log.Println(err)
		return err
	}
	defer c.Close()

	go func() {
		for err := range c.Errors() {
			log.Printf("groupId=%s, Error= %s\n", c.GroupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range c.Notifications() {
			log.Printf("groupId=%s, Rebalanced Info= %+v \n", c.GroupId, ntf)
		}
	}()

	for {
		select {
		case msg, ok := <-c.Recv():
			if ok {
				k.Ans <- msg
				//fmt.Fprintf(os.Stdout, "groupId=%s, topic=%s, partion=%d, offset=%d, key=%s, value=%s\n", c.GroupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}
	}
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
				case  <-p.Successes():
				//fmt.Println("发送成功")
				//fmt.Println("offsetCfg:", suc.Offset, " partitions:", suc.Partition, " metadata:", suc.Metadata, " value:", value)
				case fail := <-p.Errors():
					fmt.Println("err: ", fail.Err)
			}
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

//// UnSubscribe 取消订阅 consumer
//func (k *KafkaMsgImpl) UnSubscribe() error {
//	return nil
//}

//
//func (k *KafkaMsgImpl) GenerateConsumerGroup(groupname string, groupUUId string) error {
//	return nil
//}
//
//func (k *KafkaMsgImpl) JoinConsumerGroup(groupname string) error {
//	return nil
//}
//
//func (k *KafkaMsgImpl) DestroyConsumerGroup(groupname string) error {
//	return nil
//}
