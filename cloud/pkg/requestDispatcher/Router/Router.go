package Router

import (
	"fmt"
	"keep/cloud/pkg/common/kafka"

	"keep/pkg/util/core/model"
)

var RevMsgChan = make(chan *model.Message)

func MessageRouter() {

	p := kafka.NewProducerConfig("topic")
	a := kafka.NewProducerConfig("add")

	go func() { p.Publish() }()
	go func() { a.Publish() }()

	// 监听通道 路由转发
	for message := range RevMsgChan {

		switch message.Router.Group {
		case "/log":
			fmt.Println("send to fafafaf --------------------------")
			kafkaMsg := message.Content.(string)
			p.Msg <- kafkaMsg
		case "/add":
			kafkaMsg := message.Content.(string)
			a.Msg <- kafkaMsg
		}

	}

	close(p.Msg)
	close(a.Msg)
	close(RevMsgChan)
	//group := message.Router.Group
	//beehiveContext.SendToGroup(group, *message)
	// Router.MessageDispatcher(message)

}
