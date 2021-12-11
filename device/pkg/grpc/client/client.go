package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	localservice "keep/device/pkg/grpc/client/service"
	"log"
)

func main() {

	//creds,err := credentials.NewClientTLSFromFile("ca/server.crt","keep.io")
	//if err != nil {
	//	log.Fatal(err)
	//}

	client, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	rpcclient := localservice.NewProdServiceClient(client)
	response, err := rpcclient.GetProdStock(context.Background(), &localservice.ProdRequest{ProdId: 111})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.ProdStock)
}
