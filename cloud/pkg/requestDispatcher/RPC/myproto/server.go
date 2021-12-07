package coupon

import (
	context "context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	hubconfig "keep/cloud/pkg/requestDispatcher/config"
	"log"
	"net"

	certutil "k8s.io/client-go/util/cert"

	logger "keep/pkg/util/loggerv1.0.1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (this *server) SayHello(ctx context.Context, in *HelloReq) (out *HelloRsp, err error) {
	return &HelloRsp{Msg: "hello"}, nil
}

func (this *server) SayName(ctx context.Context, in *NameReq) (out *NameRsp, err error) {
	return &NameRsp{Msg: in.Name + "it is name"}, nil
}

// func EncodeCertPEM(cert *x509.Certificate) []byte {
// 	block := pem.Block{
// 		Type:  certutil.CertificateBlockType,
// 		Bytes: cert.Raw,
// 	}
// 	return pem.EncodeToMemory(&block)
// }
func ServerInit() error {
	logger.Info("gRPC server init")

	// 读取并解析公钥私钥对
	// cert, err := tls.X509KeyPair(hubconfig.Config.Cert, hubconfig.Config.Key)
	cert, err := tls.X509KeyPair(pem.EncodeToMemory(&pem.Block{Type: certutil.CertificateBlockType, Bytes: hubconfig.Config.Cert}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: hubconfig.Config.Key}))
	if err != nil {
		logger.Fatal("tls.X509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	// ca, err := ioutil.ReadFile("ca.pem")
	// if err != nil {
	// 	logger.Fatal("ioutil.ReadFile err: %v", err)
	// }

	// b := EncodeCertPEM(hubconfig.Config.Ca)
	if ok := certPool.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: certutil.CertificateBlockType, Bytes: hubconfig.Config.Ca})); !ok {
		log.Fatalf("error occurred when certPool.AppendCertsFromPEM")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	// 日志
	// grpclog.SetLoggerV2(logger.GetGrpcLogger())

	s := grpc.NewServer(
		grpc.Creds(c),
	)

	RegisterHelloServerServer(s, &server{})

	reflection.Register(s)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 9503))
	if err != nil {
		logger.Fatal("fail to listen on port 9503")
	}

	go func() {
		logger.Info("gRPC server running in goroutine")

		// register into consul
		// err = consul.Register("coupon", "127.0.0.1", 9503, "127.0.0.1:8500", time.Second*10, 15)
		// if err != nil {
		// 	panic(err)
		// }
		if err := s.Serve(listener); err != nil {
			logger.Painc("gRPC server init failed,", err)
		}

	}()

	return err
}
