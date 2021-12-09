package main

import (
	"fmt"
	counter "keep/cloud/pkg/requestDispatcher/RPC/counter"
	"net"
	"time"

	"google.golang.org/grpc"
)

type CounterServer struct{}

func (cs *CounterServer) Call(stream counter.Counter_CallServer) error {
	var n int32
	n = 0
	fmt.Println("CAll is called...")
	for {
		// _, err := stream.Recv()
		// if err == io.EOF {
		// 	break
		// }

		// if err != nil {
		// 	return err
		// }
		err := stream.Send(&counter.CounterRsp{Num: n})
		if err != nil {
			return err
		}
		n++
		time.Sleep(1 * time.Second)
		// log.Printf("SayHelloChat Client Say: %v", req.Greeting)
	}
	return nil
}
func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	counter.RegisterCounterServer(grpcServer, &CounterServer{})

	grpcServer.Serve(lis)

}
