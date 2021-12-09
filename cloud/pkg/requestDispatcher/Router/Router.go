package Router

import (
	"keep/cloud/pkg/common/kafka"
	"keep/cloud/pkg/common/modules"
	beehiveContext "keep/pkg/util/core/context"
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

func TestSendtoK8sClint() {
	// 张连军：测试抄送一份到k8sclient 可注释之
	msg_zlj := model.Message{
		Header: model.MessageHeader{},
		Router: model.MessageRoute{
			Operation: "list",
			Resource:  "$uestc/keep/k8sclient/naiveengine/pods/",
		},
		Content: map[string]string{"namespace": "default"},
	}
	beehiveContext.Send(modules.K8sClientModule, msg_zlj)
	// ==================================================
}
