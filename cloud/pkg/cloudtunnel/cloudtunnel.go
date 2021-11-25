package cloudtunnel

import (
	"github.com/wonderivan/logger"
	"keep/pkg/util/core"
)

type cloudTunnel struct {
	enable bool
}

func (c *cloudTunnel) Cleanup() {
	logger.Info("cloudTunnel clean up.")
}

func newCloudTunnel(enable bool) *cloudTunnel {
	return &cloudTunnel{
		enable: enable,
	}
}

func Register() {
	core.Register(newCloudTunnel(true))
}

func (c *cloudTunnel) Name() string {
	return "cloudTunnel"
}

func (c *cloudTunnel) Group() string {
	return "cloudTunnel"
}

func (c *cloudTunnel) Start() {
	ts := newTunnelServer()
	go ts.Start()
}

func (c *cloudTunnel) Enable() bool {
	return c.enable
}
