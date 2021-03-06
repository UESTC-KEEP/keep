package healthzhub

import (
	"github.com/UESTC-KEEP/keep/pkg/util/core/model"
	"github.com/UESTC-KEEP/keep/pkg/util/kplogger"
)

type HealthzHubImpl struct {
}

func NewHealzHub() *HealthzHubImpl {
	return new(HealthzHubImpl)
}

func (eh *HealthzHubImpl) SendMsgToModule(msg model.Message, moduleName string) error {
	go func() {
		// _, err := beehiveContext.SendSync(moduleName, msg, 5*time.Second)

		// 	if err != nil {
		// 		logger.Error("healthzagent消息发送失败:", err)
		// 	} else {
		// 		//fmt.Printf(modules.EdgeTwinModule+" 响应: %v, error: %v\n", resp, err)
		// 		logger.Trace("发送插入数据请求 " + moduleName + " 至成功...")
		// 	}
		kplogger.Error("unimpletment function ", moduleName, "\t", msg)
	}()
	return nil
}
