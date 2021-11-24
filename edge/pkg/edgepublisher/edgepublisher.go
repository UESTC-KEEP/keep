package edgepublisher

import (
	"keep/edge/pkg/common/modules"
	"keep/edge/pkg/edgepublisher/bufferpooler"
	"keep/edge/pkg/edgepublisher/chanmsgqueen"
	edgepublisherconfig "keep/edge/pkg/edgepublisher/config"
	"keep/edge/pkg/edgepublisher/publisher"
	edgeagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/edge"
	"keep/pkg/util/core"

	"github.com/wonderivan/logger"

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
	logger.Debug("准备清理模块：", modules.EdgePublisherModule)
	//bufferpooler.StopReceiveMessageForAllModules <- true
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
	go bufferpooler.StartEdgePublisher()
	bufferpooler.StartListenLogMsg()
	publisher.ReadQueueAndPublish()
}

func NewEdgePublisher(enable bool) (*EdgePublisher, error) {
	return &EdgePublisher{
		enable: enable,
	}, nil
}
