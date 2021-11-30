package kafka

import (
	"github.com/Shopify/sarama"
	saramaCluster "github.com/bsm/sarama-cluster"
	"log"
)

type ConsumerMessage = sarama.ConsumerMessage
type NotifyMessage = saramaCluster.Notification
type ConsumerError = sarama.ConsumerError

type Consumer struct {
	consumer *saramaCluster.Consumer
	Topic    string
	GroupId  string
}

func InitOneConsumerOfGroup(addr []string, topic string, groupId string, conf *Config) (*Consumer, error) {

	c, err := saramaCluster.NewConsumer(addr, groupId, []string{topic}, &conf.Config)
	var cs = &Consumer{
		consumer: c,
		Topic:    topic,
		GroupId:  groupId,
	}
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func (cs *Consumer) Close() error {
	return cs.consumer.Close()
}

func (cs *Consumer) Recv() <-chan *ConsumerMessage {
	return cs.consumer.Messages()
}

func (cs *Consumer) Notifications() <-chan *NotifyMessage {		// 发生 consumer rebalance 时通知
	return cs.consumer.Notifications()
}

func (cs *Consumer) Errors() <-chan error {
	return cs.consumer.Errors()
}

func Subscribe(address []string, topic string, groupId string ,ans chan *ConsumerMessage) {
	config :=NewConfig()
	// ans := make(chan *sarama.ConsumerMessage)
	c , err := InitOneConsumerOfGroup(address, topic, groupId, config)
	if err!=nil{
		log.Println(err)
		return
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

				ans <- msg
				//fmt.Fprintf(os.Stdout, "groupId=%s, topic=%s, partion=%d, offset=%d, key=%s, value=%s\n", c.GroupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				//c.MarkOffset(msg, "") // mark message as processed
			}
		}
	}
}

