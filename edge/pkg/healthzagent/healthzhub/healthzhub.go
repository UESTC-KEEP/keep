package healthzhub

import (
	"fmt"
	"github.com/wonderivan/logger"
	"keep/edge/pkg/common/modules"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	"time"
)

type EdgeHub struct {
}

func NewHealzHub() *EdgeHub {
	return new(EdgeHub)
}

func (eh *EdgeHub) InsertIntoSqlite(blob []byte) error {
	messsage := model.NewMessage("")
	messsage.Content = blob
	go func(blob []byte) {
		_, err := beehiveContext.SendSync(modules.EdgeTwinModule, *messsage, 5*time.Second)
		if err != nil {
			logger.Error(err)
		} else {
			//fmt.Printf(modules.EdgeTwinModule+" 响应: %v, error: %v\n", resp, err)
			fmt.Println("发送插入数据请求 " + modules.EdgeTwinModule + " 至成功...")
		}
	}(blob)
	return nil
}
