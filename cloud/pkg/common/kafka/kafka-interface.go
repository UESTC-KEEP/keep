package kafka

import "time"

// KafkaMsg kafka发送的消息格式   可以自行增删  如果kafka客户端有自己的则用之
type KafkaMsg struct {
	content  []byte
	datetime time.Time
	QoS      int
}

type KafkaInterface interface {
	// Subscribe  订阅kafka
	/*
		输入参数：host:kafka集群服务地址，为空就用默认地址
			    topic: 需要订阅的kafka主题
		返回值：表征订阅行为结果
	*/
	Subscribe(host, topic string) error
	// UnSubscribe 取消订阅
	/*
		输入参数：host：kafka集群服务地址，为空就用默认地址
				topic:要取消订阅的topic
		返回值：表征取消订阅行为的结果
	*/
	UnSubscribe(host, topic string) error
	// Publish 向kafka集群扔消息
	/*
		输入参数：msg:传递的消息结构
				topic：消息主题
				Qos:期望的服务质量
		返回值：表征发布结果
	*/
	Publish(msg *KafkaMsg, topic string, Qos int) error
	// GenerateConsumerGroup 新建订阅者组
	/*
		输入参数：groupname:生成订阅者组名
				groupUUId:指定uuid  不指定使用默认
		返回值：表征生成结果
	*/
	GenerateConsumerGroup(groupname string, groupUUId string) error
	// JoinConsumerGroup 加入消费者组  并订阅组内消息
	/*
		输入参数：groupname：需要加入的组名
	*/
	JoinConsumerGroup(groupname string) error
	// DestroyConsumerGroup 销毁消费者组
	/*
		输入参数：groupname:需要销毁的消费者组名称
	*/
	DestroyConsumerGroup(groupname string) error
}
