package cloudtunnel

import (
	"keep/cloud/pkg/requestDispatcher"
)

var Server *tunnelServer

func StartWebsocketServer() {
	ok := <-requestDispatcher.DoneTLSTunnelCerts
	if ok {
		Server = newTunnelServer()
		Server.Start()
	}
}
