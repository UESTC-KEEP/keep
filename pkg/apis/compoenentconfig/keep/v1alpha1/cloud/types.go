package cloud

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

type Modules struct {
	K8sClient           *K8sClient           `json:"k8s_client"`
	PromServer          *PromServer          `json:"prom_server"`
	RequestDispatcher   *RequestDispatcher   `json:"request_dispatcher"`
	CloudImageManager   *CloudImageManager   `json:"cloud_image_manager"`
	EqualNodeController *EqualNodeController `json:"equal_node_controller"`
}

type CloudAgentConfig struct {
	// cloudagentagent中的模块配置
	Modules *Modules `json:"modules,omitempty"`
}

// KubeAPIConfig 对看k8s apiserver进行配置
type KubeAPIConfig struct {
	// Master master的ip地址 (会把 KubeConfig 从kube文件读取的内容覆盖掉)
	// such as https://127.0.0.1:8443
	// default ""
	// Note: Can not use "omitempty" option,  It will affect the output of the default configuration file
	Master string `json:"master"`
	// ContentType 定义与k8s交互时使用的数据格式
	// default "application/vnd.kubernetes.protobuf"
	ContentType string `json:"contentType,omitempty"`
	// QPS 定义与 kubernetes apiserve 交互的并发度
	// default 100
	QPS int32 `json:"qps,omitempty"`
	// Burst 定义kubernetes apiserver交互的爆发量
	// default 200
	Burst int32 `json:"burst,omitempty"`
	// KubeConfig indicates the path to kubeConfig file with authorization and master location information.
	// default "/root/.kube/config"
	// +Required
	KubeConfig string `json:"kubeConfig"`
}

type K8sClient struct {
	Enable              bool               `json:"enable"`
	Masters             []string           `json:"masters"`
	MasterMetricTimeout int                `json:"master_metric_timeout"`
	MasterLBIp          string             `json:"master_lb_ip"`
	MasterLBPort        int                `json:"master_lb_port"`
	RedisIp             string             `json:"redis_ip"`
	RedisPort           int                `json:"redis_port"`
	PodInfo             *v1.Pod            `json:"pod_info"`
	DeploymentInfo      *appsv1.Deployment `json:"deployment_info"`
	KubeAPIConfig       *KubeAPIConfig     `json:"kubeAPIConfig,omitempty"`
	DecoderBufferSize   int                `json:"decoder_buffer_size"`
}

type PromServer struct {
	Enable                   bool `json:"enable"`
	PromServerPrometheusPort int  `json:"prom_server_prometheus_port"`
}

type RequestDispatcher struct {
	Enable        bool `json:"enable"`
	HTTPPort      int  `json:"http_port"`
	WebSocketPort int  `json:"web_socket_port"`
}

type CloudImageManager struct {
	Enable bool `json:"enable"`
}

type EqualNodeController struct {
	Enable          bool                       `json:"enable"`
	NodeName        string                     `json:"node_name"`
	MasterURL       string                     `json:"master_url"`
	KubeConfig      string                     `json:"kube_config"`
	AlsoLogToStdErr bool                       `json:"also_log_to_std_err"`
	Buffer          *EqualNodeControllerBuffer `json:"buffer"`
}

// EqualNodeControllerBuffer 定义EqualNodeController的各类通道、buffer的大小
type EqualNodeControllerBuffer struct {
	EqualNodeEvent int32 `json:"equal_node_event,omitempty"`
}
