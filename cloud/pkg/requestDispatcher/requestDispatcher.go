package requestDispatcher

import (
	"fmt"
	"keep/cloud/pkg/requestDispatcher/receiver"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"keep/pkg/util/core"
	"os"

	"github.com/wonderivan/logger"
	"k8s.io/klog/v2"
)

type RequestDispatcher struct {
	enable bool
}

func Register(r *cloudagent.RequestDispatcher) {
	rd, err := NewRequestDispatcher(r.enbale)
	if err != nil {
		logger.Error("初始化RequestDispatcher失败...:", err)
		os.Exit(1)
	}
	core.Register(rd)
}

func NewRequestDispatcher(enable bool) (*RequestDispatcher, error) {
	return &RequestDispatcher{
		enable: enable,
	}, nil
}

func (r *RequestDispatcher) Name() string {
	return ""
}
func (r *RequestDispatcher) Group() string {
	return ""
}
func (r *RequestDispatcher) Start() {

	fmt.Println("begin..")

	// check whether the certificates exist in the local directory,
	// and then check whether certificates exist in the secret, generate if they don't exist
	if err := receiver.PrepareAllCerts(); err != nil {
		klog.Exit(err)
	}
	//TODO: Will improve in the future
	//DoneTLSTunnelCerts <- true
	//close(DoneTLSTunnelCerts)

	// if err := receiver.GenerateToken(); err != nil {
	// 	klog.Exit(err)
	// }

	go receiver.StartHTTPServer()
	// HttpServer mainly used to issue certificates for the edge
	// receiver.StartReceiver()

	receiver.StartWebsocketServer()

	fmt.Println("start....")
}
func (r *RequestDispatcher) Enable() bool {
	return true
}
func (r *RequestDispatcher) Cleanup() {

}

// func main() {
// 	Start()
// }
