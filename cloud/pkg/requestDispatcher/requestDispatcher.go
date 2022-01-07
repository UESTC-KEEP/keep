package requestDispatcher

import (
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/modules"
	"github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router"
	"github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/cloudtunnel"
	requestDispatcherconfig "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/config"
	"github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/receiver"
	"github.com/UESTC-KEEP/keep/constants/cloud"
	cloudagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"github.com/UESTC-KEEP/keep/pkg/util/core"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"os"
	"time"
)

type RequestDispatcher struct {
	enable        bool
	HTTPPort      int
	WebSocketPort int
}

func Register(r *cloudagent.RequestDispatcher) {
	requestDispatcherconfig.InitConfigure(r)
	rd, err := NewRequestDispatcher(r.Enable)
	if err != nil {
		logger.Error("初始化RequestDispatcher失败...:", err)
		os.Exit(1)
	}
	core.Register(rd)
}

func NewRequestDispatcher(enable bool) (*RequestDispatcher, error) {
	return &RequestDispatcher{
		enable:        enable,
		HTTPPort:      cloud.DefaultHTTPPort,
		WebSocketPort: cloud.DefaultWebSocketPort,
	}, nil
}

func (r *RequestDispatcher) Name() string {
	return modules.RequestDispatcherModule
}
func (r *RequestDispatcher) Group() string {
	return modules.RequestDispatcherGroup
}
func (r *RequestDispatcher) Start() {

	logger.Info("RequestDispatcher begin..")

	// check whether the certificates exist in the local directory,
	// and then check whether certificates exist in the secret, generate if they don't exist
	if err := receiver.PrepareAllCerts(); err != nil {
		logger.Error(err)
	}
	//TODO: Will improve in the future
	cloudtunnel.DoneTLSTunnelCerts <- true
	close(cloudtunnel.DoneTLSTunnelCerts)

	if err := receiver.GenerateToken(); err != nil {
		logger.Fatal("fail to create Token, err:", err)
	}

	go receiver.StartHTTPServer()
	// HttpServer mainly used to issue certificates for the edge
	//receiver.StartReceiver()
	go Router.MessageRouter()

	// 测试边端发送函数
	// Router.TestRouter_SendToEdge()

	// 张连军测试
	go func() {
		time.Sleep(5 * time.Second)
		Router.TestSendtoK8sClint()
	}()

	cloudtunnel.StartWebsocketServer()

	// err := coupon.ServerInit()
	// if err != nil {
	// 	logger.Fatal("init gRPC server failed", err)
	// }
	// // 初始化路由
	// routers.InitRouters()

}
func (r *RequestDispatcher) Enable() bool {
	return r.enable
}
func (r *RequestDispatcher) Cleanup() {

}

// func main() {
// 	Start()
// }
