package k8sclient

import (
	"fmt"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/modules"
	kubeedge_device_lister "github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/kubeedge-engine/devices/lister"
	"github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/naive-engine/pods"
	edgemodules "github.com/UESTC-KEEP/keep/edge/pkg/common/modules"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/core/model"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
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
	// kubeedge设备相关
	case "$uestc/keep/k8sclient/kubeedgeengin/devices/":
		switch msg.Router.Operation {
		// 例举所有设备
		case "list":
			devicelist, _ := kubeedge_device_lister.GetAllDevice()
			logger.Trace("获取设备列表查询成功：", devicelist)
			beehiveContext.Send(modules.RequestDispatcherModule, model.Message{
				Router: model.MessageRoute{
					Source:    modules.K8sClientModule,
					Operation: "publish",
					Resource:  "edge03-zlj",
					Group:     edgemodules.DeviceMapperInterfaceGroup,
				},
				Content: devicelist,
			})
		}
	}
}
