package k8sclient

import (
	"fmt"
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/k8sclient/naive-engine/pods"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	logger "keep/pkg/util/loggerv1.0.1"
)

func StartK8sClientRouter() {
	fmt.Println("启动路由...")
	go func() {
		for {
			select {
			case <-beehiveContext.Done():
				return
			default:
			}
			msg := ReceiveMsg()
			fmt.Printf("content:%v     header:%v      router:%v\n", msg.GetContent(), msg.Header, msg.Router)
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
			listPods, err := pods.NewPods().ListPods(msg.Content.(map[string]string)["namespace"])
			if err != nil {
				logger.Error(err)
				return
			}
			fmt.Println(listPods)
		}
	}
}
