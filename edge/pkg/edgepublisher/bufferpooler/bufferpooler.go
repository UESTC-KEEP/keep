package bufferpooler

import (
	"encoding/json"
	"fmt"
	beehiveContext "keep/pkg/util/core/context"
	//beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/wonderivan/logger"
	"keep/constants"
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
var PermissionOfSending = true

// StartListenLogMsg 发送日志到消息队列中
func StartListenLogMsg() {
	go func() {
		for {
			select {
			case <-StopReceiveMessageForAllModules:
				// 收到信息停止接收所有消息
				logger.Debug("收到退出信息，清理通道...")
				PermissionOfSending = false
				return
			default:
				ReceiveFromBehiveAndPublish()

			}
		}
	}()
}

// ReceiveFromBehiveAndPublish 接收来自behivee的通信  同时返回响应 之后发布到消息队列
func ReceiveFromBehiveAndPublish() {
	msg, err := beehiveContext.Receive(modules.EdgePublisherModule)
	if err != nil {
		logger.Error(err)
		time.Sleep(5 * time.Second)
	} else {
		fmt.Printf("接收消息 msg: %v\n", msg)
		resp := msg.NewRespByMessage(&msg, " message received ")
		beehiveContext.SendResp(*resp)

		topic := constants.DefaultLogsTopic
		fmt.Println(chanmsgqueen.EdgePublishQueens)
		cli := chanmsgqueen.EdgePublishQueens[topic]
		err = cli.Publish(topic, msg.Content)
		if err != nil {
			logger.Error(err)
		}
	}
}
