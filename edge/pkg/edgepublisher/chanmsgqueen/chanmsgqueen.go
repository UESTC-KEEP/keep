package chanmsgqueen

import (
	"github.com/UESTC-KEEP/keep/constants/edge"
	"github.com/UESTC-KEEP/keep/edge/pkg/edgepublisher/config"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
)

var EdgePublishQueens = make(map[string]*Client)

// InitMsgQueens 初始化系统中需要的消息队列
func InitMsgQueens() {
	for i := 0; i < len(config.Config.EdgeMsgQueens); i++ {
		if config.Config.EdgeMsgQueens[i] == edge.DefaultLogsTopic {
			b := NewClient()
			b.SetConditions(edge.DefaultLogsQueenSize)
			EdgePublishQueens[edge.DefaultLogsTopic] = b
		} else if config.Config.EdgeMsgQueens[i] == edge.DefaultDataTopic {
			b := NewClient()
			b.SetConditions(edge.DefaultDataQueenSize)
			EdgePublishQueens[edge.DefaultDataTopic] = b
		}
	}
	//logger.Error(EdgePublishQueens)
	logger.Debug("开始初始化队列...")
	logger.Debug("队列启动成功...")
}
