package kafka

import (
	saramaCluster "github.com/bsm/sarama-cluster"
)

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

