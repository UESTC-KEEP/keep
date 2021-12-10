package Router

import (
	"keep/cloud/pkg/common/kafka"
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/requestDispatcher/Router/routers"
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
			// 指明调用函数后  功能模块返回结果的接收模块(查询pod列表后由RequestDispatcher 下发节点)
			Source: modules.RequestDispatcherModule,
			// 如果需要群发模块则填写此之段
			Group: "",
			// 以下两个内容由调用的资源模块进行解析 先定位到操作资源 在定位增删查改
			// 对资源进行的操作
			Operation: routers.KeepRouter.K8sClientRouter.NaiveEngine.Pods.Operation.List.Operation,
			// 资源所在路由
			Resource: routers.KeepRouter.K8sClientRouter.NaiveEngine.Pods.Operation.List.Resource,
		},
		// 内容及参数由RequestDispatcher与被调用模块协商
		Content: map[string]string{"namespace": "default"},
	}
	beehiveContext.Send(modules.K8sClientModule, msg_zlj)
	// ==================================================
}
