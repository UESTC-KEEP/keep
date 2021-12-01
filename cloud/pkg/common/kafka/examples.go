//ayncProducer
/*
	producer测试程序
 */

package kafka

//var address = []string{"192.168.1.103:9092", "192.168.1.103:9093"}
//func main() {
//// producer examples
//	msgValue := make(chan string , 10)
//	topic := "topic001"
//
//	go kafka.AsyncPro(address,topic,msgValue)
//
//	for i := 0;;i++ {
//		time11 := time.Now()
//		value := "this is a message 0805 "+time11.Format("15:04:05")
//		msgValue <- value
//	}
//
//	select {}
//
//// consumer examples
//	topic := "topic001"
//	ans := make(chan *sarama.ConsumerMessage , 100)
//	groupid := "cg2"
//	go kafka.Subscribe(address,topic,groupid,ans)
//
//	for msg := range ans{
//		fmt.Fprintf(os.Stdout, "groupId=%s, topic=%s, partion=%d, offset=%d, key=%s, value=%s\n", groupid, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
//	}
//	select {}
//}

