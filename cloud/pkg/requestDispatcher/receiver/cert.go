package receiver

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	hubconfig "keep/cloud/pkg/requestDispatcher/config"
	"keep/constants"

	certutil "k8s.io/client-go/util/cert"

	"keep/pkg/util/kplogger"
	"math"
	"math/big"
	"net"
	"time"
)

const validalityPeriod time.Duration = 365 * 100

func PrepareAllCerts() error {
	// Check whether the ca exists in the local directory
	if hubconfig.Config.Ca == nil && hubconfig.Config.CaKey == nil {
		kplogger.Info("Ca and CaKey creating...")
		caDER, caKey, err := NewCertificateAuthorityDer()
		if err != nil {
			kplogger.Errorf("failed to create Certificate Authority, error: %v", err)
			return err
		}

		caKeyDER, err := x509.MarshalECPrivateKey(caKey.(*ecdsa.PrivateKey))
		if err != nil {
			kplogger.Errorf("failed to convert an EC private key to SEC 1, ASN.1 DER form, error: %v", err)
			return err
		}

		UpdateConfig(caDER, caKeyDER, nil, nil)

	}

	if hubconfig.Config.Key == nil && hubconfig.Config.Cert == nil {
		kplogger.Infof("CloudCoreCert and key creating...")

		certDER, keyDER, err := SignCerts()
		if err != nil {
			kplogger.Errorf("failed to sign a certificate, error: %v", err)
			return err
		}
		UpdateConfig(nil, nil, certDER, keyDER)

	}

	return nil
}

func NewCertificateAuthorityDer() ([]byte, crypto.Signer, error) {
	caKey, err := NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	certDER, err := NewSelfSignedCACertDERBytes(caKey)
	if err != nil {
		return nil, nil, err
	}
	return certDER, caKey, nil
}

// NewPrivateKey creates an ECDSA private key
func NewPrivateKey() (crypto.Signer, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// NewSelfSignedCACertDERBytes creates a CA certificate
func NewSelfSignedCACertDERBytes(key crypto.Signer) ([]byte, error) {
	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(1024),
		Subject: pkix.Name{
			CommonName: "KeepEdge",
		},
		NotBefore: time.Now().UTC(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365 * 100),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caDERBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		return nil, err
	}
	return caDERBytes, nil
}

func UpdateConfig(ca, caKey, cert, key []byte) {
	if ca != nil {
		hubconfig.Config.Ca = ca
		kplogger.Info("update ca...")
	}
	if caKey != nil {
		hubconfig.Config.CaKey = caKey
		kplogger.Info("update caKey...")

	}
	if cert != nil {
		hubconfig.Config.Cert = cert
		kplogger.Info("update cert...")

	}
	if key != nil {
		hubconfig.Config.Key = key
		kplogger.Info("update key...")

	}
}

func getIps(advertiseAddress []string) (Ips []net.IP) {
	for _, addr := range advertiseAddress {
		Ips = append(Ips, net.ParseIP(addr))
	}
	return
}

func SignCerts() ([]byte, []byte, error) {
	cfg := &certutil.Config{
		CommonName:   "KeepEdge",
		Organization: []string{"KeepEdge"},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		AltNames: certutil.AltNames{
			DNSNames: []string{""},
			// DNSNames: hubconfig.Config.DNSNames,
			IPs: getIps([]string{constants.DefaultMasterLBIp}),
		},
	}

	certDER, keyDER, err := NewCloudCoreCertDERandKey(cfg)
	if err != nil {
		return nil, nil, err
	}

	return certDER, keyDER, nil
}

func NewCloudCoreCertDERandKey(cfg *certutil.Config) ([]byte, []byte, error) {
	serverKey, err := NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate a privateKey, err: %v", err)
	}

	keyDER, err := x509.MarshalECPrivateKey(serverKey.(*ecdsa.PrivateKey))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert an EC private key to SEC 1, ASN.1 DER form, err: %v", err)
	}

	// get ca from config
	ca := hubconfig.Config.Ca
	caCert, err := x509.ParseCertificate(ca)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse a caCert from the given ASN.1 DER data, err: %v", err)
	}

	caKeyDER := hubconfig.Config.CaKey
	caKey, err := x509.ParseECPrivateKey(caKeyDER)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse ECPrivateKey, err: %v", err)
	}

	certDER, err := NewCertFromCa(cfg, caCert, serverKey.Public(), caKey, validalityPeriod)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate a certificate using the given CA certificate and key, err: %v", err)
	}
	return certDER, keyDER, nil
}
func NewCertFromCa(cfg *certutil.Config, caCert *x509.Certificate, serverKey crypto.PublicKey, caKey crypto.Signer, validalityPeriod time.Duration) ([]byte, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	if len(cfg.CommonName) == 0 {
		return nil, errors.New("must specify a CommonName")
	}
	if len(cfg.Usages) == 0 {
		return nil, errors.New("must specify at least one ExtKeyUsage")
	}

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:     cfg.AltNames.DNSNames,
		IPAddresses:  cfg.AltNames.IPs,
		SerialNumber: serial,
		NotBefore:    time.Now().UTC(),
		NotAfter:     time.Now().Add(time.Hour * 24 * validalityPeriod),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  cfg.Usages,
	}
	certDERBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, caCert, serverKey, caKey)
	if err != nil {
		return nil, err
	}
	return certDERBytes, nil
}
