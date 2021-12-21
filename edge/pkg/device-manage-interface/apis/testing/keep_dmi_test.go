package device_manage_interface

import (
	"context"
	"fmt"
	localservice "keep/edge/pkg/device-manage-interface/apis/devices/v1"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestDMI(t *testing.T) {

	//creds,err := credentials.NewClientTLSFromFile("ca/server.crt","keep.io")
	//if err != nil {
	//	log.Fatal(err)
	//}

	go TestDummyMapper(t)

	client, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	dev := localservice.NewDeviceServiceClient(client)
	for i := 0; i < 4; i++ {
		response, err := dev.GetDeviceStatus(context.Background(), &localservice.DeviceStatusRequest{Topic: "test_topic"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(response.String())
		time.Sleep(time.Second)
	}

}
