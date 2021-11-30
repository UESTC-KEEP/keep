package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

//var Address = []string{"192.168.1.103:9092", "192.168.1.103:9093"}
type ProducerMessage = sarama.ProducerMessage	// bie ming
type ProducerError = sarama.ProducerError

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
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return nil, err
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

func AsyncPro(address []string, topic string, msg  chan string)  {
	config := NewConfig()
	//config.Producer.Return.Successes = true
	//config.Producer.Return.Errors = true
	produ,_ := InitManualRetryAsyncProducer(address, config)
	defer produ.Close()
	go func (p *AsyncProducer){
		for{
			select {
			case suc := <-p.Successes():
				//case  <-p.Successes():
				//fmt.Println("发送成功")
				//bytes, _ := suc.Value.Encode()
				//value := string(bytes)
				//fmt.Println("offsetCfg:", suc.Offset, " partitions:", suc.Partition, " metadata:", suc.Metadata, " value:", value)
			case fail := <-p.Errors():
				fmt.Println("err: ", fail.Err)
			}
		}
	}(produ)

	for msgvalue := range msg{

		//发送的消息,主题,key
		msg := &ProducerMessage{
			Topic: topic,
		}

		//将字符串转化为字节数组
		msg.Value = sarama.ByteEncoder(msgvalue)

		//使用通道发送
		produ.Send() <- msg

		time.Sleep(500*time.Millisecond)
	}
}
