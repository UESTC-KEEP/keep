package bufferpooler

import (
	"encoding/json"
	"github.com/wonderivan/logger"
	"keep/constants"
	"keep/edge/pkg/edgepublisher/chanmsgqueen"
	"keep/edge/pkg/edgepublisher/config"
	"net/http"
	"strconv"
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

// SendLogInQueue 发送日志到消息队列中
func SendLogInQueue(log string) {
	topic := constants.DefaultLogsTopic
	cli := chanmsgqueen.EdgePublishQueens[topic]
	err := cli.Publish(topic, log)
	if err != nil {
		logger.Error(err)
	}
}
