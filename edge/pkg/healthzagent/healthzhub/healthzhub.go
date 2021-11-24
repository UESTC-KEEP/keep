package healthzhub

import (
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	"keep/pkg/util/kplogger"
	"time"
)

type HealthzHubImpl struct {
}

func NewHealzHub() *HealthzHubImpl {
	return new(HealthzHubImpl)
}

func (eh *HealthzHubImpl) SendMsgToModule(msg model.Message, moduleName string) error {
	go func() {
		_, err := beehiveContext.SendSync(moduleName, msg, 5*time.Second)
		if err != nil {
			kplogger.Error(err)
		} else {
			//fmt.Printf(modules.EdgeTwinModule+" 响应: %v, error: %v\n", resp, err)
			kplogger.Trace("发送插入数据请求 " + moduleName + " 至成功...")
		}
	}()
	return nil
}
