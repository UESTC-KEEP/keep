package edgepublisher

import (
	"encoding/json"
	"fmt"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgepublisher/bufferpooler"
	"keep/edge/pkg/edgepublisher/chanmsgqueen"
	edgepublisherconfig "keep/edge/pkg/edgepublisher/config"
	"keep/edge/pkg/edgepublisher/publisher"
	edgetunnel "keep/edge/pkg/edgepublisher/tunnel"
	"keep/edge/pkg/edgepublisher/tunnel/cert"
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"keep/pkg/util/core"
	logger "keep/pkg/util/loggerv1.0.1"
	"net/http"
	"strconv"

	"os"
)

type EdgePublisher struct {
	enable            bool
	httpserver        string
	port              int32
	heartbeat         int32
	tlscafile         string
	tlscertfile       string
	tlsprivatekeyfile string
	edgemsgqueens     []string
	hostnameOverride  string
	nodeIP            string
	token             string
}

// Register 注册healthzagent模块
func Register(ep *edgeagent.EdgePublisher) {
	edgepublisherconfig.InitConfigure(ep)
	edgepublisher, err := NewEdgePublisher(ep.Enable)
	if err != nil {
		logger.Error("初始化logsagent失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(edgepublisher)
}

func (ep *EdgePublisher) Cleanup() {

}

func (ep *EdgePublisher) Name() string {
	return modules.EdgePublisherModule
}

func (ep *EdgePublisher) Group() string {
	return modules.EdgePublisherGroup
}

//Enable indicates whether this module is enabled
func (ep *EdgePublisher) Enable() bool {
	return ep.enable
}

func (ep *EdgePublisher) Start() {
	//var wg sync.WaitGroup
	name, _ := os.Hostname()
	fmt.Println("name:", name)
	logger.Debug("EdgePublisher 开始启动....")
	// 申请证书
	nm := cert.NewCertManager("NodeName", ep.token)
	nm.Start()
	go edgetunnel.StartEdgeTunnel(ep.hostnameOverride, ep.nodeIP)

	// 启动边端服务20350
	// 初始化队列 确保队列初始化完成再启动服务
	chanmsgqueen.InitMsgQueens()
	//wg.Wait()
	go StartEdgePublisher()
	go bufferpooler.StartListenLogMsg()
	go publisher.ReadQueueAndPublish()

	//err := coupon.CouponClientInit()
	//if err != nil {
	//	logger.Fatal("init coupon gRPC client failed", err)
	//}
}

// StartEdgePublisher 边端健康检测 20350端口的/用于云端对边端进行健康性  存活性检测
func StartEdgePublisher() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", EdgeAgentHealthCheck)
	logger.Debug("edgepublisher  :" + strconv.Itoa(int(edgepublisherconfig.Config.ServePort)) + " 服务启动中...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(edgepublisherconfig.Config.ServePort)), mux)
	if err != nil {
		logger.Fatal("publisher启动失败,端口占用：", err)
	}
}

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

func NewEdgePublisher(enable bool) (*EdgePublisher, error) {
	return &EdgePublisher{
		enable:           enable,
		hostnameOverride: edgepublisherconfig.Config.HostnameOverride,
		nodeIP:           edgepublisherconfig.Config.LocalIP,
		token:            "630f3e6b363262ab595d14eb9665b48b40c2b70f591244d6dc32cd43bb69e83d.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1Mjk4Njd9.5zuMK9YIv9BWlv2_QypgRQysIGpWFni4LbxhRHJ36ys",
	}, nil
}
