package k8sclient

import (
	"fmt"
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/k8sclient/naive-engine/pods"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	logger "keep/pkg/util/loggerv1.0.1"
)

// SendBeehiveMsg 发送beehive消息
func SendBeehiveMsg(module string, msg model.Message) {
	beehiveContext.Send(module, msg)
}

func StartK8sClientRouter() {
	//fmt.Println("启动路由...")
	go func() {
		for {
			select {
			case <-beehiveContext.Done():
				return
			default:
			}
			msg := ReceiveMsg()
			fmt.Printf("%#v\n", msg)
			if msg != nil {
				// 来自其他模块的普通消息
				if msg.Router.Resource == "" {

				} else {
					// 来自router的调用
					ResolveRouter(msg)
				}
			}
		}
	}()
}

func ReceiveMsg() *model.Message {
	msg, err := beehiveContext.Receive(modules.K8sClientModule)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return &msg
}

func ResolveRouter(msg *model.Message) {
	switch msg.Router.Resource {
	case "$uestc/keep/k8sclient/naiveengine/pods/":
		switch msg.Router.Operation {
		case "list":
			msgMap := msg.Content.(map[string]interface{})

			listPods, err := pods.NewPods().ListPods(msgMap["namespace"].(string))
			if err != nil {
				logger.Error(err)
				return
			}
			fmt.Print("pod 列表：")
			for _, pod := range listPods.Items {
				fmt.Print(pod.Name + "   ")
			}
			//SendBeehiveMsg()
			//fmt.Println(listPods)
		}
	}
}
