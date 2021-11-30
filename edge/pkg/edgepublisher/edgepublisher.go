package edgepublisher

import (
	"encoding/json"
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgepublisher/bufferpooler"
	"keep/edge/pkg/edgepublisher/chanmsgqueen"
	edgepublisherconfig "keep/edge/pkg/edgepublisher/config"
	"keep/edge/pkg/edgepublisher/publisher"
	edgetunnel "keep/edge/pkg/edgepublisher/tunnel"
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"keep/pkg/util/core"
	"keep/pkg/util/loggerv1.0.0"
	"net/http"
	"strconv"

	"os"
	"sync"
)

// func main() {
// 	nm := publisher.NewCertManager("NodeName")
// 	nm.Start()

// 	pool := x509.NewCertPool()
// 	caCertPath := "/etc/kubeedge/ca/rootCA.crt"

// 	caCrt, err := ioutil.ReadFile(caCertPath)
// 	if err != nil {
// 		fmt.Println("ReadFile err:", err)
// 		return
// 	}
// 	pool.AppendCertsFromPEM(caCrt)

// 	cliCrt, err := tls.LoadX509KeyPair(constants.DefaultCertFile, constants.DefaultKeyFile)
// 	if err != nil {
// 		fmt.Println("Loadx509keypair err:", err)
// 		return
// 	}

// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{
// 			RootCAs:      pool,
// 			Certificates: []tls.Certificate{cliCrt},
// 		},
// 	}
// 	client := &http.Client{Transport: tr}
// 	//这里的ip地址需要在生成自签名证书的时候指定,否则ssl验证不通过。
// 	res, err := client.Get("https://192.168.1.121:2022")
// 	if err != nil {
// 		fmt.Println("client get error:", err)
// 	}
// 	defer res.Body.Close()
// 	body, _ := ioutil.ReadAll(res.Body)
// 	fmt.Println(string(body))

// }

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
}

// Register 注册healthzagent模块
func Register(ep *edgeagent.EdgePublisher, nodeName, nodeIP string) {
	edgepublisherconfig.InitConfigure(ep)
	edgepublisher, err := NewEdgePublisher(ep.Enable, nodeName, nodeIP)
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

func (l *EdgePublisher) Start() {
	var wg sync.WaitGroup
	logger.Debug("EdgePublisher 开始启动....")
	// 启动边端服务20350
	// 初始化队列 确保队列初始化完成再启动服务
	chanmsgqueen.InitMsgQueens()
	wg.Wait()
	go StartEdgePublisher()
	bufferpooler.StartListenLogMsg()
	publisher.ReadQueueAndPublish()
	go edgetunnel.StartEdgeTunnel(l.hostnameOverride, l.nodeIP)
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

func NewEdgePublisher(enable bool, hostnameOverride, nodeIP string) (*EdgePublisher, error) {
	return &EdgePublisher{
		enable:           enable,
		hostnameOverride: hostnameOverride,
		nodeIP:           nodeIP,
	}, nil
}
