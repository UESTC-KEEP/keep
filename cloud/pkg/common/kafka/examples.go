//ayncProducer
/*
	producer测试程序
 */

package kafka

import (
	"kafka/kafka"
	"time"
)
var address = []string{"192.168.1.103:9092", "192.168.1.103:9093"}
func main() {
	msgValue := make(chan string , 10)
	topic := "topic001"

	go kafka.AsyncPro(address,topic,msgValue)

	for i := 0;;i++ {
		time11 := time.Now()
		value := "this is a message 0805 "+time11.Format("15:04:05")
		msgValue <- value
	}

	select {}
}

