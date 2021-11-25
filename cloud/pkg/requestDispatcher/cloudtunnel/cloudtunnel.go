package cloudtunnel

func StartWebsocketServer() {
	ts := newTunnelServer()
	go ts.Start()
}
