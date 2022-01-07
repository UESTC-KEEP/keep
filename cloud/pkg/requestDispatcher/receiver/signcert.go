package receiver

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	hubconfig "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/config"
	"github.com/UESTC-KEEP/keep/constants/cloud"

	certutil "k8s.io/client-go/util/cert"

	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"net"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func SignCerts() ([]byte, []byte, error) {
	cfg := &certutil.Config{
		CommonName:   "KeepEdge",
		Organization: []string{"KeepEdge"},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		AltNames: certutil.AltNames{
			DNSNames: []string{""},
			// DNSNames: hubconfig.Config.DNSNames,
			IPs: getIps([]string{cloud.DefaultKeepCloudIP}),
		},
	}

	certDER, keyDER, err := NewCloudCoreCertDERandKey(cfg)
	if err != nil {
		return nil, nil, err
	}

	return certDER, keyDER, nil
}

func getIps(advertiseAddress []string) (Ips []net.IP) {
	for _, addr := range advertiseAddress {
		Ips = append(Ips, net.ParseIP(addr))
	}
	return
}

// GenerateToken will create a token consisting of caHash and jwt Token and save it to secret
func GenerateToken() error {
	// set double TokenRefreshDuration as expirationTime, which can guarantee that the validity period
	// of the token obtained at anytime is greater than or equal to TokenRefreshDuration
	expiresAt := time.Now().Add(time.Hour * hubconfig.Config.RequestDispatcher.TokenRefreshDuration * 2).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.StandardClaims{
		ExpiresAt: expiresAt,
	}

	keyPEM := getCaKey()
	tokenString, err := token.SignedString(keyPEM)

	if err != nil {
		return fmt.Errorf("failed to generate the token for EdgeCore register, err: %v", err)
	}

	caHash := getCaHash()
	// combine caHash and tokenString into caHashAndToken
	caHashToken := strings.Join([]string{caHash, tokenString}, ".")
	// save caHashAndToken to secret
	// err = CreateTokenSecret([]byte(caHashToken))
	// if err != nil {
	// 	return fmt.Errorf("failed to create tokenSecret, err: %v", err)
	// }
	UpdateToken(caHashToken)

	t := time.NewTicker(time.Hour * hubconfig.Config.RequestDispatcher.TokenRefreshDuration)
	go func() {
		for {
			<-t.C
			refreshedCaHashToken := refreshToken()
			// if err := CreateTokenSecret([]byte(refreshedCaHashToken)); err != nil {
			// 	klog.Exitf("failed to create the ca token for edgecore register, err: %v", err)
			// }
			UpdateToken(refreshedCaHashToken)
		}
	}()
	logger.Info("Succeed to creating token")
	return nil
}

func refreshToken() string {
	claims := &jwt.StandardClaims{}
	expirationTime := time.Now().Add(time.Hour * hubconfig.Config.RequestDispatcher.TokenRefreshDuration * 2)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	keyPEM := getCaKey()
	tokenString, err := token.SignedString(keyPEM)
	if err != nil {
		logger.Error("Failed to generate token signed by caKey, err: %v", err)
	}
	caHash := getCaHash()
	//put caHash in token
	caHashAndToken := strings.Join([]string{caHash, tokenString}, ".")
	return caHashAndToken
}

// getCaHash gets ca-hash
func getCaHash() string {
	caDER := hubconfig.Config.Ca
	digest := sha256.Sum256(caDER)
	return hex.EncodeToString(digest[:])
}

// getCaKey gets caKey to encrypt token
func getCaKey() []byte {
	caKey := hubconfig.Config.CaKey
	return caKey
}
