package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"keep/device/pkg/grpc/services"
	"log"
	"net/http"
)

func main() {
	flag.Parse()
	rpcServer := grpc.NewServer()
	services.RegisterProdServiceServer(rpcServer, new(services.ProductService))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.Proto)
		fmt.Println(request.Header)
		rpcServer.ServeHTTP(writer, request)
	})

	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	err := httpServer.ListenAndServeTLS("ca/server.crt", "ca/server_no_passwd.key")
	if err != nil {
		log.Fatal(err)
	}

}
