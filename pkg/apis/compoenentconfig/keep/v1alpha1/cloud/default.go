package cloud

import (
	eqndconstants "keep/cloud/pkg/equalnodecontroller/constants"
	"keep/constants"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

// NewDefaultEdgeAgentConfig returns a full EdgeCoreConfig object
func NewDefaultEdgeAgentConfig() *CloudAgentConfig {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		// 如果输入了kubeconfig参数，该参数的值就是kubeconfig文件的绝对路径，
		// 如果没有输入kubeconfig参数，就用默认路径~/.kube/config
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		// 如果取不到当前用户的家目录，就没办法设置kubeconfig的默认目录了，只能从入参中取
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	return &CloudAgentConfig{
		Modules: &Modules{
			K8sClient: &K8sClient{
				Enable:              true,
				Masters:             strings.Split(constants.DefaultMasterIpPort, ";"),
				MasterMetricTimeout: constants.DefaultMasterMetricTimeout,
				MasterLBIp:          constants.DefaultMasterLBIp,
				MasterLBPort:        constants.DefaultMasterLBPort,
				RedisIp:             constants.DefaultRedisServerIp,
				RedisPort:           constants.DefaultRedisServerPort,
				PodInfo:             nil,
				DeploymentInfo:      nil,
				KubeAPIConfig: &KubeAPIConfig{
					//Master:      "",
					//ContentType: "",
					QPS:        100,
					Burst:      200,
					KubeConfig: *kubeconfig,
				},
				DecoderBufferSize: constants.DefaultDecoderBufferSize,
			},
			PromServer: &PromServer{
				Enable:                   true,
				PromServerPrometheusPort: constants.DefaultPromServerMetricsPort,
			},
			RequestDispatcher: &RequestDispatcher{
				Enable:               true,
				HTTPPort:             constants.DefaultHTTPPort,
				WebSocketPort:        constants.DefaultWebSocketPort,
				TokenRefreshDuration: constants.DefaultTokenRefreshDuration,
			},
			CloudImageManager: &CloudImageManager{
				Enable: true,
			},
			EqualNodeController: &EqualNodeController{
				Enable:          true,
				MasterURL:       constants.DefaultMasterURL,
				KubeConfig:      *kubeconfig,
				AlsoLogToStdErr: constants.DefaultAlsoLogToStdErr,
				Buffer: &EqualNodeControllerBuffer{
					EqualNodeEvent: eqndconstants.DefaultEqualNodeEventBuffer,
				},
			},
		},
	}
}
