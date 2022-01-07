package receiver

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/UESTC-KEEP/keep/constants/cloud"
	"io/ioutil"
	"net/http"
	"strings"

	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog"

	hubconfig "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/config"
)

const (
	Address = "0.0.0.0"
	Port    = cloud.DefaultHTTPPort
)

// StartHTTPServer starts the http service
func StartHTTPServer() {
	// gorilla/mux是 gorilla Web 开发工具包中的路由管理库
	router := mux.NewRouter()
	//  HandleFunc 返回一个路由  Methods为路由添加接受的请求类型：GET, POST
	router.HandleFunc(cloud.DefaultCertURL, edgeCoreClientCert).Methods(http.MethodGet)
	router.HandleFunc(cloud.DefaultCAURL, getCA).Methods(http.MethodGet)

	addr := fmt.Sprintf("%s:%d", Address, Port)
	// 509 公钥认证的标准格式  tls安全传输协议 pem数据格式
	// pem 实现了 PEM 数据编码，起源 Privacy Enhanced Mail。今天 PEM 编码最常见的用途是在 TLS 密钥和证书中
	cert, err := tls.X509KeyPair(pem.EncodeToMemory(&pem.Block{Type: certutil.CertificateBlockType, Bytes: hubconfig.Config.Cert}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: hubconfig.Config.Key}))

	if err != nil {
		logger.Fatal(err)
	}
	// 创建一个http server
	server := &http.Server{
		Addr:    addr,
		Handler: router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequestClientCert,
		},
	}
	// TLS
	logger.Info("listening....")
	logger.Fatal(server.ListenAndServeTLS("", ""))
}

// getCA returns the caCertDER
func getCA(w http.ResponseWriter, r *http.Request) {
	logger.Info("getting ca...")
	caCertDER := hubconfig.Config.Ca
	if _, err := w.Write(caCertDER); err != nil {
		logger.Error("failed to write caCertDER, err: %v", err)
	}
}

// func edgeCert(w http.ResponseWriter, r *http.Request) {
// 	klog.Info("edge call....")
// 	certDER := hubconfig.Config.Cert
// 	keyDER := hubconfig.Config.Key

// }

// EncodeCertPEM returns PEM-encoded certificate data
func EncodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  certutil.CertificateBlockType,
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

// edgeCoreClientCert will verify the certificate of EdgeCore or token then create EdgeCoreCert and return it
func edgeCoreClientCert(w http.ResponseWriter, r *http.Request) {
	logger.Info("getting cert...")

	// 用于tls验证  http用户忽略
	if cert := r.TLS.PeerCertificates; len(cert) > 0 {
		if err := verifyCert(cert[0]); err != nil {
			logger.Error("failed to sign the certificate for edgenode: %s, failed to verify the certificate", r.Header.Get(cloud.NodeName))
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				logger.Error("failed to write response, err: %v", err)
			}
		} else {
			signEdgeCert(w, r)
		}
		return
	}
	// 因为是http 直接验证权限
	if verifyAuthorization(w, r) {
		signEdgeCert(w, r)
	} else {
		logger.Error("failed to sign the certificate for edgenode: %s, invalid token", r.Header.Get(cloud.NodeName))
	}
}

// verifyCert verifies the edge certificate by CA certificate when edge certificates rotate.
func verifyCert(cert *x509.Certificate) error {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: certutil.CertificateBlockType, Bytes: hubconfig.Config.Ca}))
	if !ok {
		return fmt.Errorf("failed to parse root certificate")
	}
	opts := x509.VerifyOptions{
		Roots:     roots,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify edge certificate: %v", err)
	}
	return nil
}

// verifyAuthorization verifies the token from EdgeCore CSR
func verifyAuthorization(w http.ResponseWriter, r *http.Request) bool {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
			klog.Error("failed to write http response, err: %v", err)
		}
		return false
	}
	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
			klog.Error("failed to write http response, err: %v", err)
		}
		return false
	}
	// 验证token的合法性
	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		caKey := hubconfig.Config.CaKey
		return caKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
				logger.Error("Write body error %v", err)
			}
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
			logger.Error("Write body error %v", err)
		}

		return false
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
			klog.Error("Write body error %v", err)
		}
		return false
	}
	return true
}

// signEdgeCert signs the CSR from EdgeCore
func signEdgeCert(w http.ResponseWriter, r *http.Request) {
	csrContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("fail to read file when signing the cert for edgenode:%s! error:%v", r.Header.Get(cloud.NodeName), err)
		return
	}
	csr, err := x509.ParseCertificateRequest(csrContent)
	if err != nil {
		logger.Error("fail to ParseCertificateRequest of edgenode: %s! error:%v", r.Header.Get(cloud.NodeName), err)
		return
	}
	usagesStr := r.Header.Get("ExtKeyUsages")
	var usages []x509.ExtKeyUsage
	if usagesStr == "" {
		usages = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	} else {
		err := json.Unmarshal([]byte(usagesStr), &usages)
		if err != nil {
			logger.Error("unmarshal http header ExtKeyUsages fail, err: %v", err)
			return
		}
	}
	logger.Info("receive sign crt request, ExtKeyUsages: %v", usages)
	clientCertDER, err := signCerts(csr.Subject, csr.PublicKey, usages)
	if err != nil {
		logger.Error("fail to signCerts for edgenode:%s! error:%v", r.Header.Get(cloud.NodeName), err)
		return
	}

	if _, err := w.Write(clientCertDER); err != nil {
		logger.Error("write error %v", err)
	}
}

// signCerts will create a certificate for EdgeCore
func signCerts(subInfo pkix.Name, pbKey crypto.PublicKey, usages []x509.ExtKeyUsage) ([]byte, error) {
	cfgs := &certutil.Config{
		CommonName:   subInfo.CommonName,
		Organization: subInfo.Organization,
		Usages:       usages,
	}
	clientKey := pbKey

	ca := hubconfig.Config.Ca
	caCert, err := x509.ParseCertificate(ca)
	if err != nil {
		return nil, fmt.Errorf("unable to ParseCertificate: %v", err)
	}

	caKeyDER := hubconfig.Config.CaKey
	caKey, err := x509.ParseECPrivateKey(caKeyDER)
	if err != nil {
		return nil, fmt.Errorf("unable to ParseECPrivateKey: %v", err)
	}

	// edgeCertSigningDuration := hubconfig.Config.CloudHub.EdgeCertSigningDuration
	certDER, err := NewCertFromCa(cfgs, caCert, clientKey, caKey, 365) //crypto.Signer(caKey)
	if err != nil {
		return nil, fmt.Errorf("unable to NewCertFromCa: %v", err)
	}

	return certDER, err
}
