package edgetunnel

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/UESTC-KEEP/keep/constants/cloud"
	"github.com/UESTC-KEEP/keep/constants/edge"
	"github.com/UESTC-KEEP/keep/edge/pkg/common/modules"
	"github.com/UESTC-KEEP/keep/pkg/stream"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/core/model"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type edgeTunnel struct {
	hostnameOverride string
	nodeIP           string
}

var session *tunnelSession                                                //封装后的websocket连接
var sessionConnected bool                                                 //云边是否连通
var msgSendBuffer = make([]*model.Message, edge.DefaultMsgSendBufferSize) //发送到云的缓冲
var msgSendBufferLock sync.Locker                                         //缓冲锁
var reconnectChan = make(chan struct{})                                   //需要重连时向此channel发

func newEdgeTunnel(hostnameOverride, nodeIP string) *edgeTunnel {
	return &edgeTunnel{
		hostnameOverride: hostnameOverride,
		nodeIP:           nodeIP,
	}
}

func (e *edgeTunnel) start() {
	serverURL := url.URL{
		Scheme: "wss",
		Host:   fmt.Sprintf("%s:%d", edge.DefaultHttpServer, edge.DefaultWebSocketPort),
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

		go session.startPing(reconnectChan)
		go routeToEdge()
		//go routeToCloud()

		<-reconnectChan
		sessionConnected = false
		session.Close()
		logger.Warn("connection broken, reconnecting...")
		time.Sleep(5 * time.Second)

		//清空reconnectChan
	clean:
		for {
			select {
			case <-reconnectChan:
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
	var err error
	if session != nil {
		err = session.Tunnel.WriteMessage([]*model.Message{msg})
	} else {
		err = errors.New("websocket session is nil,try write to edgetwin...")
		logger.Trace("send message to edge cloud error: ", err)
	}

	if err != nil {
		_, err = beehiveContext.SendSync(modules.EdgeTwinModule, *msg, time.Second)
		if err != nil {
			logger.Warn("send message to edge twin error: ", err)
		}
	}
}

func WriteToBufferToCloud(msg *model.Message) {
	var err error
	msgSendBufferLock.Lock()
	msgSendBuffer = append(msgSendBuffer, msg)
	if len(msgSendBuffer) == edge.DefaultMsgSendBufferSize {
		if session != nil {
			err = session.Tunnel.WriteMessage(msgSendBuffer)
		} else {
			err = errors.New("websocket session is nil")
		}

		if err != nil {
			for _, contentMsg := range msgSendBuffer {
				_, err := beehiveContext.SendSync(modules.EdgeTwinGroup, *contentMsg, time.Second)
				if err != nil {
					logger.Error("send message to edge twin error: ", err)
				}
			}
		}
		msgSendBuffer = msgSendBuffer[0:0]
	}
	msgSendBufferLock.Unlock()
}

func routeToCloud() {
	for {
		select {
		case <-beehiveContext.Done():
			logger.Warn("EdgeTunnel RouteToEdge stop")
			return
		default:
		}

		message, err := beehiveContext.Receive(modules.EdgePublisherModule)
		if err != nil {
			logger.Error("failed to receive message from edge: ", err)
			time.Sleep(time.Second)
			continue
		}

		WriteToCloud(&message)

	}
}

func routeToEdge() {
	for {
		select {
		case <-beehiveContext.Done():
			logger.Warn("EdgeTunnel RouteToEdge stop")
			return
		default:
		}

		_, r, err := session.Tunnel.NextReader()
		if err != nil {
			logger.Error("Read messsage error: ", err)
			reconnectChan <- struct{}{}
			return
		}

		messList, err := stream.ReadMessageFromTunnel(r)
		if err != nil {
			logger.Error("Get tunnel message error: ", err)
			reconnectChan <- struct{}{}
			return
		}

		//如果是对某条消息的响应消息
		for _, contentMsg := range messList {
			if contentMsg.Header.ParentID != "" {
				beehiveContext.SendResp(*contentMsg)
			} else {
				fmt.Printf("%#v", contentMsg)
				beehiveContext.SendToGroup(contentMsg.Router.Group, *contentMsg)
			}
		}

	}
}
