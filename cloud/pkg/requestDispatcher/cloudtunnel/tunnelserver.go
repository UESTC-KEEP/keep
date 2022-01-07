package cloudtunnel

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/modules"
	"github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/config"
	"github.com/UESTC-KEEP/keep/constants/cloud"
	"github.com/UESTC-KEEP/keep/pkg/stream"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/core/model"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/gorilla/websocket"
)

type tunnelServer struct {
	container  *restful.Container
	upgrader   websocket.Upgrader
	sync.Mutex //
	sessions   map[string]*session
	nodeNameIP sync.Map
}

func newTunnelServer() *tunnelServer {
	return &tunnelServer{
		container: restful.NewContainer(),
		sessions:  make(map[string]*session),
		upgrader: websocket.Upgrader{
			HandshakeTimeout: time.Second * 2,
			ReadBufferSize:   1024,
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				w.WriteHeader(status)
				_, err := w.Write([]byte(reason.Error()))
				logger.Error("failed to write http response, err: ", err)
			},
		},
	}
}

func (s *tunnelServer) installDefaultHandler() {
	ws := new(restful.WebService)
	ws.Path(cloud.DefaultWebSocketUrl)
	ws.Route(ws.GET("/").
		To(s.connect))
	s.container.Add(ws)
}

func (s *tunnelServer) addsession(key string, session *session) {
	s.Lock()
	s.sessions[key] = session
	s.Unlock()
}

func (s *tunnelServer) getSession(id string) (*session, bool) {
	s.Lock()
	defer s.Unlock()
	sess, ok := s.sessions[id]
	return sess, ok
}

func (s *tunnelServer) addNodeIP(node, ip string) {
	s.nodeNameIP.Store(node, ip)
}

func (s *tunnelServer) getNodeIP(node string) (string, bool) {
	ip, ok := s.nodeNameIP.Load(node)
	if !ok {
		return "", ok
	}
	return ip.(string), ok
}

func (s *tunnelServer) connect(r *restful.Request, w *restful.Response) {
	hostnameOverride := r.HeaderParameter(cloud.SessionKeyHostNameOverride)
	internalIP := r.HeaderParameter(cloud.SessionKeyInternalIP)
	if internalIP == "" {
		internalIP = strings.Split(r.Request.RemoteAddr, ":")[0]
	}
	con, err := s.upgrader.Upgrade(w, r.Request, nil)
	if err != nil {
		return
	}
	logger.Info("get a new tunnel agent: ", hostnameOverride, internalIP)
	session := &session{
		tunnel:    stream.NewDefaultTunnel(con),
		sessionID: hostnameOverride,
	}

	s.addsession(hostnameOverride, session)
	s.addsession(internalIP, session)
	s.addNodeIP(hostnameOverride, internalIP)
	session.Serve()
}

func (s *tunnelServer) Start() {
	s.installDefaultHandler()

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: config.Config.Ca}))

	certificate, err := tls.X509KeyPair(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: config.Config.Cert}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: config.Config.Key}))
	if err != nil {
		logger.Error("Failed to load TLSTunnelCert and key")
		panic(err)
	}

	tunnelServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cloud.DefaultWebSocketPort),
		Handler: s.container,
		TLSConfig: &tls.Config{
			ClientCAs:    pool,
			Certificates: []tls.Certificate{certificate},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			MinVersion:   tls.VersionTLS12,
			CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
		},
	}
	logger.Info("Prepare to start tunnel server")

	go s.sendBeehiveMessageToEdge()

	err = tunnelServer.ListenAndServeTLS("", "")
	if err != nil {
		logger.Fatal("Start tunnelServer error", err)
		return
	}
}

func (s *tunnelServer) sendBeehiveMessageToEdge() {
	logger.Info("send message to edge goroutine started")
	for {
		select {
		case <-beehiveContext.Done():
			return
		default:
		}
		msg, err := beehiveContext.Receive(modules.RequestDispatcherModule)
		fmt.Printf("%#v\n", msg)
		logger.Info("send message to edge: ", msg)
		if err != nil {
			logger.Info("receive not Message format message")
			continue
		}

		nodeID := msg.Router.Resource
		if nodeID == "" {
			logger.Warn("node id not found in the message")
			continue
		}

		session, ok := s.getSession(nodeID)
		if !ok {
			logger.Error("node: ", nodeID, "doesn't exist, send message error")
			continue
		}
		err = session.writeMessageToTunnel(&msg)
		if err != nil {
			logger.Error("write to tunnel error, edgenode", nodeID)
			continue
		}
	}
}

func (s *tunnelServer) SendToEdge(edge string, msg *model.Message) {
	session, ok := s.getSession(edge)
	if !ok {
		logger.Error("no such edge node")
	}

	err := session.writeMessageToTunnel(msg)
	if err != nil {
		logger.Error("write tunnel error: ", err)
	}
}
