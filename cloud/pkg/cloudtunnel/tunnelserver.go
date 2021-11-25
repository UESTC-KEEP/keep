package cloudtunnel

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"keep/pkg/stream"
	beehiveContext "keep/pkg/util/core/context"
	"net/http"
	"strings"
	"sync"
	"time"
)

type TunnelServer struct {
	container  *restful.Container
	upgrader   websocket.Upgrader
	sync.Mutex //
	sessions   map[string]*Session
	nodeNameIP sync.Map
}

func newTunnelServer() *TunnelServer {
	return &TunnelServer{
		container: restful.NewContainer(),
		sessions:  make(map[string]*Session),
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

func (s *TunnelServer) installDefaultHandler() {
	ws := new(restful.WebService)
	ws.Path("/v1/keepedge/connect")
	ws.Route(ws.GET("/").
		To(s.connect))
	s.container.Add(ws)
}

func (s *TunnelServer) addSession(key string, session *Session) {
	s.Lock()
	s.sessions[key] = session
	s.Unlock()
}

func (s *TunnelServer) getSession(id string) (*Session, bool) {
	s.Lock()
	defer s.Unlock()
	sess, ok := s.sessions[id]
	return sess, ok
}

func (s *TunnelServer) addNodeIP(node, ip string) {
	s.nodeNameIP.Store(node, ip)
}

func (s *TunnelServer) getNodeIP(node string) (string, bool) {
	ip, ok := s.nodeNameIP.Load(node)
	if !ok {
		return "", ok
	}
	return ip.(string), ok
}

func (s *TunnelServer) connect(r *restful.Request, w *restful.Response) {
	hostnameOverride := r.HeaderParameter("SessionHostNameOverride")
	internalIP := r.HeaderParameter("SessionInternalIP")
	if internalIP == "" {
		internalIP = strings.Split(r.Request.RemoteAddr, ":")[0]
	}
	con, err := s.upgrader.Upgrade(w, r.Request, nil)
	if err != nil {
		return
	}
	logger.Info("get a new tunnel agent: ", hostnameOverride, internalIP)
	session := &Session{
		tunnel:    stream.NewDefaultTunnel(con),
		sessionID: hostnameOverride,
	}

	s.addSession(hostnameOverride, session)
	s.addSession(internalIP, session)
	s.addNodeIP(hostnameOverride, internalIP)
	session.Serve()
}

func (s *TunnelServer) Start() {
	s.installDefaultHandler()
	var data []byte
	var key []byte
	var cert []byte

	data, err := ioutil.ReadFile("/etc/keepedge/ca/rootCA.crt")
	if err == nil {
		block, _ := pem.Decode(data)
		data = block.Bytes
	}

	key, err = ioutil.ReadFile("/etc/keepedge/certs/server.key")
	if err == nil {
		block, _ := pem.Decode(key)
		key = block.Bytes
	}

	cert, err = ioutil.ReadFile("/etc/keepedge/certs/server.crt")
	if err == nil {
		block, _ := pem.Decode(cert)
		cert = block.Bytes
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: data}))

	certificate, err := tls.X509KeyPair(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key}))
	if err != nil {
		logger.Error("Failed to load TLSTunnelCert and key")
		panic(err)
	}

	tunnelServer := &http.Server{
		Addr:    fmt.Sprintf("%d", 3721),
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
	err = tunnelServer.ListenAndServeTLS("", "")
	if err != nil {
		logger.Fatal("Start tunnelServer error", err)
		return
	}
	go s.sendMessageToEdge()
}

func (s *TunnelServer) sendMessageToEdge() {
	for {
		msg, err := beehiveContext.Receive("cloudTunnel")
		logger.Info("[cloudTunnel] send message to edge: ", msg)
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
		err = session.WriteMessageToTunnel(&msg)
		if err != nil {
			logger.Error("[cloudtunnel] write to tunnel error, edgenode", nodeID)
			continue
		}
	}
}
