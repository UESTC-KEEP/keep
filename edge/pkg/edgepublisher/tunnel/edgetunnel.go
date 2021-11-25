package edgetunnel

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	beehiveContext "keep/pkg/util/core/context"
	"net/http"
	"net/url"
	"time"
)

type edgeTunnel struct {
	hostnameOverride string
	nodeIP           string
	reconnectChan    chan struct{}
}

func newEdgeTunnel(hostnameOverride, nodeIP string) *edgeTunnel {
	return &edgeTunnel{
		hostnameOverride: hostnameOverride,
		nodeIP:           nodeIP,
		reconnectChan:    make(chan struct{}),
	}
}

func (e *edgeTunnel) Start() {
	serverURL := url.URL{
		Scheme: "wss",
		Host:   "192.168.1.140:3721",
		Path:   "/v1/keepedge/connect",
	}

	cert, err := tls.LoadX509KeyPair("/etc/kubeedge/certs/server.crt", "/etc/kubeedge/certs/server.key")
	if err != nil {
		logger.Fatal("Failed to load x509 key pair: ", err)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	for {
		select {
		case <-beehiveContext.Done():
			return
		default:
		}
		session, err := e.TLSClientConnect(serverURL, tlsConfig)
		if err != nil {
			logger.Error("connect failed: ", err)
			time.Sleep(5 * time.Second)
			continue
		}

		go session.startPing(e.reconnectChan)
		go session.routeToEdge(e.reconnectChan)

		<-e.reconnectChan
		session.Close()
		logger.Warn("connection broken, reconnecting...")
		time.Sleep(5 * time.Second)

		//清空reconnectChan
	clean:
		for {
			select {
			case <-e.reconnectChan:
			default:
				break clean
			}
		}
	}
}

func (e *edgeTunnel) TLSClientConnect(url url.URL, tlsConfig *tls.Config) (*TunnelSession, error) {
	logger.Info("Start a new tunnel connection")

	dial := websocket.Dialer{
		TLSClientConfig:  tlsConfig,
		HandshakeTimeout: time.Duration(30) * time.Second,
	}
	header := http.Header{}
	header.Add("SessionHostNameOverride", e.hostnameOverride)
	header.Add("SessionInternalIP", e.nodeIP)

	con, _, err := dial.Dial(url.String(), header)
	if err != nil {
		return nil, err
	}

	session := NewTunnelSession(con)
	return session, nil
}

func (e *edgeTunnel) keepAlive() {

}
