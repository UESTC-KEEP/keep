package conf

import (
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)
var KubeMaster=""
var Kubeconfig=""
//var Kubeconfig=filepath.Join(
//	os.Getenv("HOME"),".kube","config")
var KubeQPS = float32(5.000000)
var KubeBurst = 10
var KubeContentType = "application/vnd.kubernetes.protobuf"
func GetKubeConfig() (*rest.Config,error) {
	if home:=homedir.HomeDir();home!="" {
		Kubeconfig=flag.String("kubeconfig",filepath.Join(home,".kube","config"))
	}
	config, err := clientcmd.BuildConfigFromFlags(KubeMaster, Kubeconfig)
	config.QPS=KubeQPS
	config.Burst=KubeBurst
	config.ContentType=KubeContentType
	return config,err
}
