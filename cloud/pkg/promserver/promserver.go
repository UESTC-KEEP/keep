package promserver

import (
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/modules"
	"github.com/UESTC-KEEP/keep/cloud/pkg/promserver/config"
	cloudagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"github.com/UESTC-KEEP/keep/pkg/util/core"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"os"
)

type PromServer struct {
	enable bool
}

func Register(ps *cloudagent.PromServer) {
	config.InitConfigure(ps)
	promserver, err := NewPromServer(ps.Enable)
	if err != nil {
		logger.Error("初始化promserver失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(promserver)
}

func (ps *PromServer) Cleanup() {}

func (ps *PromServer) Start() {

}

func (ps *PromServer) Name() string {
	return modules.PromServerModule
}
func (ps *PromServer) Group() string {
	return modules.PromServerGroup
}

func (ps *PromServer) Enable() bool {
	return ps.enable
}

func NewPromServer(enable bool) (*PromServer, error) {
	return &PromServer{
		enable: enable,
	}, nil
}
