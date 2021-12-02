package edgetunnel

import (
	"crypto/tls"
	"fmt"
	"keep/constants"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	logger "keep/pkg/util/loggerv1.0.1"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type edgeTunnel struct {
	hostnameOverride string
	nodeIP           string
	reconnectChan    chan struct{}
}

var session *tunnelSession
var sessionConnected bool

func newEdgeTunnel(hostnameOverride, nodeIP string) *edgeTunnel {
	return &edgeTunnel{
		hostnameOverride: hostnameOverride,
		nodeIP:           nodeIP,
		reconnectChan:    make(chan struct{}),
	}
}

func (e *edgeTunnel) start() {
	serverURL := url.URL{
		Scheme: "wss",
		Host:   fmt.Sprintf("%s:%d", constants.DefaultKeepCloudIP, constants.DefaultWebSocketPort),
		Path:   constants.DefaultWebSocketUrl,
	}

	cert, err := tls.LoadX509KeyPair(constants.DefaultCertFile, constants.DefaultKeyFile)
	if err != nil {
		logger.Info("Failed to load x509 key pair: ", err, "try again")
		time.Sleep(10 * time.Second)
		cert, err = tls.LoadX509KeyPair(constants.DefaultCertFile, constants.DefaultKeyFile)
	}
	if err != nil {
		logger.Fatal("Failed to load x509 key pair: ", err, "Exiting...")
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
		var err error
		session, err = e.tlsClientConnect(serverURL, tlsConfig)
		if err != nil {
			logger.Error("connect failed: ", err)
			time.Sleep(5 * time.Second)
			continue
		}
		sessionConnected = true

		go session.startPing(e.reconnectChan)
		go session.routeToEdge(e.reconnectChan)

		<-e.reconnectChan
		sessionConnected = false
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

func (e *edgeTunnel) tlsClientConnect(url url.URL, tlsConfig *tls.Config) (*tunnelSession, error) {
	logger.Info("Start a new tunnel connection")

	dial := websocket.Dialer{
		TLSClientConfig:  tlsConfig,
		HandshakeTimeout: time.Duration(30) * time.Second,
	}
	header := http.Header{}
	header.Add(constants.SessionKeyHostNameOverride, e.hostnameOverride)
	header.Add(constants.SessionKeyInternalIP, e.nodeIP)

	con, _, err := dial.Dial(url.String(), header)
	if err != nil {
		return nil, err
	}

	session := NewTunnelSession(con)
	return session, nil
}

func StartEdgeTunnel(nodeName, nodeIP string) {
	edget := newEdgeTunnel(nodeName, nodeIP)
	edget.start()
}

func WriteToCloud(msg *model.Message) error {
	for !sessionConnected {
		logger.Info("session not connected, waiting")
		time.Sleep(3 * time.Second)
	}
	return session.Tunnel.WriteMessage(msg)

}
