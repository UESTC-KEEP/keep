package main

import (
	"flag"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"keep/device/pkg/grpc/services"
	"net"
)

func main() {
	flag.Parse()

	//creds,err := credentials.NewServerTLSFromFile("ca/server.crt","ca/server_no_passwd.key")
	//if err != nil {
	//	glog.Errorf("%#v",err)
	//}

	rpcServer := grpc.NewServer()
	services.RegisterProdServiceServer(rpcServer, new(services.ProductService))
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		glog.Errorf("%#v", err)
	}
	err = rpcServer.Serve(listener)
	if err != nil {
		glog.Errorf("%#v", err)
	}
}
