package LogWatcher

import (
	"bytes"
	"encoding/json"
	"github.com/wonderivan/logger"
	"keep/cloud/pkg/common/kafka"
	"keep/constants"
	"net/http"
)

type LogStruct struct {
	Logid string `json:"logid"`

}

func GetLogFromKafka(){
	err := kafka.KafkaInterface.Subscribe("topic")
	if err!=nil{
	logger.Error(err.Error())
	}
	ch:=make(chan string)
	read:=<-ch
	marshal, err := json.Marshal(read)
	buffer := bytes.NewBuffer(marshal)
	_, err = http.Post(constants.Url, constants.ContentType, buffer)
	if err!=nil{
		logger.Error(err.Error())
	}

}

