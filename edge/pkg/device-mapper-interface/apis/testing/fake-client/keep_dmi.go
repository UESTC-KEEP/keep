package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	localservice "keep/edge/pkg/device-mapper-interface/apis/devices/v1"
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

	rpcclient := localservice.NewDeviceServiceClient(client)
	response, err := rpcclient.GetDeviceStatus(context.Background(), &localservice.DeviceStatusRequest{Topic: "test_topic"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.String())
}
