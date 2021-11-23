package receiver

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	hubconfig "keep/cloud/pkg/requestDispatcher/config"
	"net/http"

	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog"
)

type myhandler struct {
}

func (h *myhandler) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprintf(w,
		"Hi, This is an example of http service in golang!\n")
}

func createTLSConfig(ca, cert, key []byte) tls.Config {
	// init certificate
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: certutil.CertificateBlockType, Bytes: ca}))
	if !ok {
		panic(fmt.Errorf("fail to load ca content"))
	}

	certificate, err := tls.X509KeyPair(pem.EncodeToMemory(&pem.Block{Type: certutil.CertificateBlockType, Bytes: cert}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key}))
	if err != nil {
		panic(err)
	}
	return tls.Config{
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS12,
		// has to match cipher used by NewPrivateKey method, currently is ECDSA
		CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
	}
}

func StartWebsocketServer() {
	tlsConfig := createTLSConfig(hubconfig.Config.Ca, hubconfig.Config.Cert, hubconfig.Config.Key)
	svc := &http.Server{
		// Type:       api.ProtocolTypeWS,
		TLSConfig: &tlsConfig,
		Handler:   &myhandler{},
		// AutoRoute:  true,
		// ConnNotify: handler.CloudhubHandler.OnRegister,
		Addr: (":20000"),
		// ExOpts:     api.WSServerOption{Path: "/"},
	}
	// klog.Infof("Starting cloudhub %s server", api.ProtocolTypeWS)
	fmt.Println("websocket listening...")
	klog.Exit(svc.ListenAndServeTLS("", ""))
	fmt.Println("stop websocket listen")
}
