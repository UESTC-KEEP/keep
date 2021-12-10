package kafka

import (
	"keep/constants/cloud"
	"log"
)

type ConsumerConfig struct {
	Address []string
	Topic   string
	GroupId string
	Ans     chan *ConsumerMessage
}

func NewConsumerConfig(topic string, groupid string) *ConsumerConfig {
	msg := make(chan *ConsumerMessage, 100)
	c := &ConsumerConfig{
		Address: cloud.Address,
		//Address: []string{"192.168.1.103:9092", "192.168.1.103:9093"},
		Topic:   topic,
		GroupId: groupid,
		Ans:     msg,
	}
	return c
}

// Subscribe 订阅/接收消息 consumer
func (k *ConsumerConfig) Subscribe() error {
	config := NewConfig()
	// ans := make(chan *sarama.ConsumerMessage)
	c, err := InitOneConsumerOfGroup(k.Address, k.Topic, k.GroupId, config)
	if err != nil {
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
