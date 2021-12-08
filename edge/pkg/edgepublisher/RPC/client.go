package coupon

import (
	context "context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"keep/constants"

	"io/ioutil"

	logger "keep/pkg/util/loggerv1.0.1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
)

func CouponClientInit() error {
	var err error

	// 读取并解析公钥私钥对
	cert, err := tls.LoadX509KeyPair(constants.DefaultCertFile, constants.DefaultKeyFile)
	if err != nil {
		logger.Fatal("Load tls.LoadX509KeyPair error: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(constants.DefaultCAFile)
	if err != nil {
		logger.Fatal("ioutil.ReadFile error: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		logger.Fatal("certPool.AppendCertsFromPEM error")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "",
		RootCAs:      certPool,
	})

	// set log
	// grpclog.SetLoggerV2(logger.GetGrpcLogger())
	target := "192.168.1.121:9503"

	// 连接到注册中心
	Conn, err := grpc.Dial(
		target,
		grpc.WithBalancerName(roundrobin.Name),
		grpc.WithTransportCredentials(c),
	)
	defer Conn.Close()
	if err != nil {
		logger.Fatal("coupon grpc failed to connect to the given target: %v", err)
	}

	//获得grpc句柄
	conn := NewHelloServerClient(Conn)
	// couponClient = NewCouponClient(couponServerConn)

	re1, err := conn.SayHello(context.Background(), &HelloReq{Name: "songguojun"})
	if err != nil {
		fmt.Println("calling SayHello() error", err)
	}
	fmt.Println(re1.Msg)

	//通过句柄进行调用服务端函数SayName
	re2, err := conn.SayName(context.Background(), &NameReq{Name: "songguojun"})
	if err != nil {
		fmt.Println("calling SayName() error", err)
	}
	fmt.Println(re2.Msg)
	return err
}
