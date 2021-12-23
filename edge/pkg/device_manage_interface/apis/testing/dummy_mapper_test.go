package device_manage_interface

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"testing"

	v1 "keep/edge/pkg/device-manage-interface/apis/devices/v1"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

//随机生成数据并传给grpc
type DevObj struct {
}

func (dev *DevObj) GetDeviceStatus(ctx context.Context, req *v1.DeviceStatusRequest) (*v1.DeviceStatusResponse, error) {
	dev_id := req.Topic
	fmt.Println("Get req=", dev_id)

	temp := strconv.Itoa(20 + rand.Intn(10))
	return &v1.DeviceStatusResponse{
		Status: &v1.DeviceStatus{

			Status: map[string]string{"temperature": temp},
			Info:   map[string]string{"result": "ok"},
		},
		Info: map[string]string{dev_id: "ok"},
	}, nil
}

func TestDummyMapper(t *testing.T) {
	flag.Parse()

	//creds,err := credentials.NewServerTLSFromFile("ca/server.crt","ca/server_no_passwd.key")
	//if err != nil {
	//	glog.Errorf("%#v",err)
	//}

	rpcServer := grpc.NewServer()
	v1.RegisterDeviceServiceServer(rpcServer, &DevObj{})
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		glog.Errorf("%#v", err)
	}
	fmt.Println("start dummy mapper")
	err = rpcServer.Serve(listener)
	if err != nil {
		glog.Errorf("%#v", err)
	}
}
