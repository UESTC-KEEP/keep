package LogWatcher

import (
	"keep/cloud/pkg/common/kafka"
	"keep/constants"
)

func GetAndPusherKafka(){
	messages := make(chan *kafka.ConsumerMessage, 100)
    go kafka.Subscribe([]string{constants.Address},constants.OrginTopic,"logstash",messages)

	var mtlog =make(chan string,100)
	go kafka.AsyncPro([]string{constants.Address},constants.ParseTopic,mtlog)

	for msg:=range messages{
		value := string(msg.Value) + "lll"
		mtlog<-value
	}


	//ms :="lnf1"


}



