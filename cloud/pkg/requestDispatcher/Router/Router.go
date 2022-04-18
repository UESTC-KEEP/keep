package Router

import (
	"fmt"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/kafka"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/modules"
	"github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router/routers"
	prome_exporter "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/prome-exporter"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/core/model"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
)

var RevMsgChan = make(chan *model.Message)

func MessageRouter() {
	routers.InitRouters()
	p := kafka.NewProducerConfig("topic")
	a := kafka.NewProducerConfig("add")

	go func() { p.Publish() }()
	go func() { a.Publish() }()

	// 监听通道 路由转发
	// for message := range RevMsgChan {
	for {
		select {
		case <-beehiveContext.Done():
			close(p.Msg)
			close(a.Msg)
			close(RevMsgChan)
			return
		case message := <-RevMsgChan:
			msgStr := fmt.Sprintf("云端接收消息:%#v", message)
			//if len(msgStr) >= 60 {
			//	msgStr = msgStr[:60]
			//}
			logger.Trace(msgStr+"..."+"message.Router.Resource:", message.Router.Resource)
			switch message.Router.Resource {
			case "/log":
				kafkaMsg := message.Content.(string)
				p.Msg <- kafkaMsg
			case "/add":
				kafkaMsg := message.Content.(string)
				a.Msg <- kafkaMsg
			// 匹配naive_engine pod资源 ToDo: 应该合并两类资源  同时传送给K8sClientModule
			case routers.KeepRouter.K8sClientRouter.NaiveEngine.Pods.Resources:
				beehiveContext.Send(modules.K8sClientModule, *message)
			// 匹配kubeedge device资源
			case routers.KeepRouter.K8sClientRouter.KubeedgeEngine.Devices.Resources:
				beehiveContext.Send(modules.K8sClientModule, *message)
			// 匹配需要发送给kafka的用量消息
			case routers.KeepRouter.Kafka.Metrics.Resources:
				logger.Info("发送给prometheus...")
				//node_name := message.Content.(map[string]interface{})["node_name"].(string)
				//cpuUsage := message.Content.(map[string]interface{})["cpu_usage"].(float64)
				//memUsage := message.Content.(map[string]interface{})["mem"].(map[string]interface{})["usedPercent"].(float64)
				//prome_exporter.UpdatesData(node_name, cpuUsage, memUsage)
				prome_exporter.UpdateNodeExporter(message.Content.(string))
			}
			// 匹配metrics信息

		default:
		}
	}

	// }

}

func TestSendtoK8sClint() {
	// 张连军：测试抄送一份到k8sclient 可注释之
	testmap := make(map[string]interface{})
	testmap["namespace"] = "default"
	msg_zlj := model.Message{
		Header: model.MessageHeader{},
		Router: model.MessageRoute{
			// 指明调用函数后  功能模块返回结果的接收模块(查询pod列表后由RequestDispatcher 下发节点)
			Source: modules.RequestDispatcherModule,
			// 如果需要群发模块则填写此之段
			Group: "",
			// 以下两个内容由调用的资源模块进行解析 先定位到操作资源 在定位增删查改
			// 对资源进行的操作
			Operation: routers.KeepRouter.K8sClientRouter.KubeedgeEngine.Devices.Operation.List,
			// 资源所在路由
			Resource: routers.KeepRouter.K8sClientRouter.KubeedgeEngine.Devices.Resources,
		},
		// 内容及参数由RequestDispatcher与被调用模块协商
		Content: testmap,
	}
	beehiveContext.Send(modules.K8sClientModule, msg_zlj)
	// ==================================================
}

var SendChan = make(chan model.Message)

func SendToEdge() {
	// for {
	// 	select {
	// 	case <-beehiveContext.Done():
	// 		close(SendChan)
	// 		return
	// 	default:
	// 	}
	// 	msg := <-SendChan

	// }
	for msg := range SendChan {
		fmt.Println("message:", msg)
	}

}

// 测试函数
func TestRouter_SendToEdge() {

	msg := &model.Message{}
	msg.Content = "hello edge!!!"
	msg.Router.Group = "/log"

	beehiveContext.AddModule("router")

	go SendToEdge()

	beehiveContext.Send("router", *msg)
	// time.Sleep(3 * time.Second)

	message := <-SendChan
	fmt.Println("------------------msg :", message.Content)

}
