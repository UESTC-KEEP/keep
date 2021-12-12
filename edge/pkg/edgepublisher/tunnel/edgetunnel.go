package edgetunnel

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"keep/constants/cloud"
	"keep/constants/edge"
	"keep/edge/pkg/common/modules"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
	logger "keep/pkg/util/loggerv1.0.1"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type edgeTunnel struct {
	hostnameOverride string
	nodeIP           string
	reconnectChan    chan struct{}
}

var session *tunnelSession
var sessionConnected bool
var msgSendBuffer = make([]*model.Message, edge.DefaultMsgSendBufferSize)
var msgSendBufferLock sync.Locker

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
		Host:   fmt.Sprintf("%s:%d", cloud.DefaultKeepCloudIP, cloud.DefaultWebSocketPort),
		Path:   cloud.DefaultWebSocketUrl,
	}
	//certManager := cert.NewCertManager(e.hostnameOverride, config.Config.Token)
	//certManager.Start()

	clientCert, err := tls.LoadX509KeyPair(edge.DefaultCertFile, edge.DefaultKeyFile)
	if err != nil {
		logger.Info("Failed to load x509 key pair: ", err, "try again")
		time.Sleep(10 * time.Second)
		clientCert, err = tls.LoadX509KeyPair(edge.DefaultCertFile, edge.DefaultKeyFile)
	}
	if err != nil {
		logger.Fatal("Failed to load x509 key pair: ", err, "Exiting...")
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{clientCert},
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
	header.Add(cloud.SessionKeyHostNameOverride, e.hostnameOverride)
	header.Add(cloud.SessionKeyInternalIP, e.nodeIP)

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

func WriteToCloud(msg *model.Message) {
	for i := 0; i < 5 && !sessionConnected; i++ {
		logger.Info("session not connected, waiting")
		time.Sleep(3 * time.Second)
	}
	if !sessionConnected {
		msgToEdgeTwin := model.NewMessage("")
		msgToEdgeTwin.SetResourceOperation(msg.GetResource(), "")
		_, err := beehiveContext.SendSync(modules.EdgeTwinGroup, *msgToEdgeTwin, time.Second)
		if err != nil {
			logger.Error("send message to edge twin error: ", err)
		}
		return
	}
	err := session.Tunnel.WriteMessage([]*model.Message{msg})
	if err != nil {
		msgToEdgeTwin := model.NewMessage("")
		msgToEdgeTwin.SetResourceOperation(msg.GetResource(), "")
		_, err := beehiveContext.SendSync(modules.EdgeTwinGroup, *msgToEdgeTwin, time.Second)
		if err != nil {
			logger.Error("send message to edge twin error: ", err)
		}
	}

}

func WriteToBufferToCloud(msg *model.Message) {
	for i := 0; i < 5 && !sessionConnected; i++ {
		logger.Info("session not connected, waiting")
		time.Sleep(3 * time.Second)
	}

	msgSendBufferLock.Lock()
	msgSendBuffer = append(msgSendBuffer, msg)
	if len(msgSendBuffer) == edge.DefaultMsgSendBufferSize {
		err := session.Tunnel.WriteMessage(msgSendBuffer)
		if err != nil {
			for _, contentMsg := range msgSendBuffer {
				msgToEdgeTwin := model.NewMessage("")
				msgToEdgeTwin.SetResourceOperation(contentMsg.GetResource(), "")
				_, err := beehiveContext.SendSync(modules.EdgeTwinGroup, *msgToEdgeTwin, time.Second)
				if err != nil {
					logger.Error("send message to edge twin error: ", err)
				}
			}
		}
		msgSendBuffer = msgSendBuffer[0:0]
	}
	msgSendBufferLock.Unlock()
}
