package main

import (
	"context"
	"fmt"
	"io"
	counter "keep/cloud/pkg/requestDispatcher/RPC/counter"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("grpc dial failed err:", err)
	}
	client := counter.NewCounterClient(conn)
	rep, err := client.Call(context.Background())
	if err != nil {
		fmt.Errorf("grpc call failed")
	}
	for {
		res, err := rep.Recv()
		if err == io.EOF {
			fmt.Println("connect break")
			break
		}
		if err != nil {
			fmt.Println("call counter err: ", err)
		}
		fmt.Println("call counter : ", res.Num)

	}

}
