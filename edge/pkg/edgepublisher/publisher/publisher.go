// Package publisher 负责对消息队列中的数据 或者  非消息队列的即时数据进行发送到云端
package publisher

import (
	"fmt"
	"keep/edge/pkg/edgepublisher/chanmsgqueen"
	"keep/edge/pkg/edgepublisher/config"
	edgetunnel "keep/edge/pkg/edgepublisher/tunnel"
	"keep/pkg/util/core/model"
	logger "keep/pkg/util/loggerv1.0.1"
)

//
func ReadQueueAndPublish() {
	// 对每一个topic都启动一个携程 监听队列
	for i := 0; i < len(config.Config.EdgeMsgQueens); i++ {
		topic := config.Config.EdgeMsgQueens[i]
		cli := chanmsgqueen.EdgePublishQueens[topic]
		//logger.Error(cli)
		ch, err := cli.Subscribe(topic)
		if err != nil {
			logger.Error(err)
		}
		go func() {
			for {
				e := cli.GetPayLoad(ch)
				Publish(e)
			}
		}()
	}
}

// Publish 将数据发送到云端
func Publish(msg interface{}) {
	fmt.Println("--------------------------  发送云端  msg:", msg)
	message := model.Message{}
	message.Content = msg
	message.Router.Group = "/log"
	edgetunnel.WriteToCloud(&message)
}
