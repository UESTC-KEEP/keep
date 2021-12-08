package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	coupon "keep/cloud/pkg/requestDispatcher/RPC/myproto"
	"net"
	"time"

	logger "keep/pkg/util/loggerv1.0.1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	// 其他省略
)

func ServerInit() error {
	logger.Info("gRPC server init")

	// 读取并解析公钥私钥对
	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		logger.Fatal("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		logger.Fatal("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		logger.Fatal("error occurred when certPool.AppendCertsFromPEM")
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

	coupon.RegisterCouponServer(s, &couponService.Service{})

	reflection.Register(s)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 9503))
	if err != nil {
		logger.Fatal("fail to listen on port 9503")
	}

	go func() {
		logger.Info("gRPC server running in goroutine")

		// register into consul
		err = consul.Register("coupon", "127.0.0.1", 9503, "127.0.0.1:8500", time.Second*10, 15)
		if err != nil {
			panic(err)
		}
		if err := s.Serve(listener); err != nil {
			logger.Painc("gRPC server init failed,", err)
		}
	}()

	return err
}
