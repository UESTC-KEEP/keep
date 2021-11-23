package main

import (
	"fmt"
	"keep/cloud/pkg/requestDispatcher/receiver"

	"k8s.io/klog/v2"
)

func Start() {

	fmt.Println("begin..")

	// check whether the certificates exist in the local directory,
	// and then check whether certificates exist in the secret, generate if they don't exist
	if err := receiver.PrepareAllCerts(); err != nil {
		klog.Exit(err)
	}
	// TODO: Will improve in the future
	//DoneTLSTunnelCerts <- true
	//close(DoneTLSTunnelCerts)

	// if err := receiver.GenerateToken(); err != nil {
	// 	klog.Exit(err)
	// }

	go receiver.StartHTTPServer()
	// HttpServer mainly used to issue certificates for the edge
	// receiver.StartReceiver()

	receiver.StartWebsocketServer()

	fmt.Println("start....")
}

func main() {
	Start()
}
