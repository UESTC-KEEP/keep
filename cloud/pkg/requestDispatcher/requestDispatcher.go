package requestDispatcher

import (
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/requestDispatcher/Router"
	"keep/cloud/pkg/requestDispatcher/Router/routers"
	"keep/cloud/pkg/requestDispatcher/cloudtunnel"
	requestDispatcherconfig "keep/cloud/pkg/requestDispatcher/config"
	"keep/cloud/pkg/requestDispatcher/receiver"
	"keep/constants/cloud"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/core"
	logger "keep/pkg/util/loggerv1.0.1"
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

	cloudtunnel.StartWebsocketServer()

	// err := coupon.ServerInit()
	// if err != nil {
	// 	logger.Fatal("init gRPC server failed", err)
	// }
	// 初始化路由
	routers.InitRouters()
	// 张连军测试
	go func() {
		time.Sleep(2 * time.Second)
		Router.TestSendtoK8sClint()
	}()
}
func (r *RequestDispatcher) Enable() bool {
	return r.enable
}
func (r *RequestDispatcher) Cleanup() {

}

// func main() {
// 	Start()
// }
