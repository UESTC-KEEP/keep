package cloudtunnel

var Server *tunnelServer
var DoneTLSTunnelCerts = make(chan bool, 1)

func StartWebsocketServer() {
	ok := <-DoneTLSTunnelCerts
	if ok {
		Server = newTunnelServer()
		Server.Start()
	}
}
