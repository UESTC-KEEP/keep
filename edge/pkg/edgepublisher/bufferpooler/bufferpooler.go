package bufferpooler

import (
	"encoding/json"
	"fmt"
	"github.com/wonderivan/logger"
	"keep/constants"
	beehiveContext "keep/core/context"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgepublisher/chanmsgqueen"
	"keep/edge/pkg/edgepublisher/config"
	"net/http"
	"strconv"
	"time"
)

// SentImmediately 被调用就即时给云端推送消息
//func SentImmediately(){
//	port := config.Config.Port
//	server := config.Config.HTTPServer
//
//}

type response struct {
	Content string `json:"content"`
	Errinfo string `json:"errinfo,omitempty"`
}

func EdgeAgentHealthCheck(w http.ResponseWriter, r *http.Request) {
	res := response{
		Content: "edgeagent工作中",
		Errinfo: "",
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&res)
	if err != nil {
		logger.Error(err)
	}
}

func InitCachePools() {

}

// StartEdgePublisher 边端健康检测
func StartEdgePublisher() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", EdgeAgentHealthCheck)
	logger.Debug("edgepublisher  :" + strconv.Itoa(int(config.Config.ServePort)) + " 服务启动中...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.Config.ServePort)), mux)
	if err != nil {
		logger.Error(err)
	}
}

var StopReceiveMessageForAllModules = make(chan bool)

// SendLogInQueue 发送日志到消息队列中
func SendLogInQueue() {
	go func() {
		for {
			select {
			case <-StopReceiveMessageForAllModules:
				logger.Debug("收到退出信息，清理通道...")
				close(StopReceiveMessageForAllModules)
				time.Sleep(3 * time.Second)
			default:
				msg, err := beehiveContext.Receive(modules.EdgePublisherModule)
				if err != nil {
					logger.Error(err)
				}
				fmt.Println("接收消息 msg: %v\n", msg)
				resp := msg.NewRespByMessage(&msg, " message received ")
				beehiveContext.SendResp(*resp)

				topic := constants.DefaultLogsTopic
				cli := chanmsgqueen.EdgePublishQueens[topic]
				err = cli.Publish(topic, msg.Content)
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}()
}
