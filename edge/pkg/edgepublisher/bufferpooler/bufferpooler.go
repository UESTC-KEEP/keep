package bufferpooler

import (
	"fmt"
	"keep/constants/edge"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/loggerv1.0.1"

	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgepublisher/chanmsgqueen"
)

// SentImmediately 被调用就即时给云端推送消息
//func SentImmediately(){
//	port := config.Config.Port
//	server := config.Config.HTTPServer
//
//}

func InitCachePools() {

}

// StartListenLogMsg 发送日志到消息队列中
func StartListenLogMsg() {
	go func() {
		for {
			select {
			case <-beehiveContext.Done():
				// 收到信息停止接收所有消息
				//PermissionOfSending = false
				return
			default:
			}
			ReceiveFromBeehiveAndPublish()
		}
	}()
}

// ReceiveFromBeehiveAndPublish 接收来自behivee的通信  同时返回响应 之后发布到消息队列
func ReceiveFromBeehiveAndPublish() {
	msg, err := beehiveContext.Receive(modules.EdgePublisherModule)
	if err != nil {
		logger.Error(err)
		return
	}
	fmt.Printf("接收消息 msg: %#v\n", msg)
	topic := edge.DefaultDataTopic
	//fmt.Println(chanmsgqueen.EdgePublishQueens)
	cli := chanmsgqueen.EdgePublishQueens[topic]
	err = cli.Publish(topic, msg)
	if err != nil {
		logger.Error(err)
	}
}
