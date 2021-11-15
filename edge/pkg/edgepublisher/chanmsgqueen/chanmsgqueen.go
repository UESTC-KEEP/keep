package chanmsgqueen

import (
	"github.com/wonderivan/logger"
	"keep/constants"
	"keep/edge/pkg/edgepublisher/config"
)

var EdgePublishQueens = make(map[string]*Client)

// InitMsgQueens 初始化系统中需要的消息队列
func InitMsgQueens() {
	for i := 0; i < len(config.Config.EdgeMsgQueens); i++ {
		if config.Config.EdgeMsgQueens[i] == constants.DefaultLogsTopic {
			b := NewClient()
			b.SetConditions(constants.DefaultLogsQueenSize)
			EdgePublishQueens[constants.DefaultLogsTopic] = b
		} else if config.Config.EdgeMsgQueens[i] == constants.DefaultDataTopic {
			b := NewClient()
			b.SetConditions(constants.DefaultDataQueenSize)
			EdgePublishQueens[constants.DefaultDataTopic] = b
		}
	}
	//logger.Error(EdgePublishQueens)
	logger.Debug("开始初始化队列...")
	logger.Debug("队列启动成功...")
}
