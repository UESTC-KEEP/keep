package cloud

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

type Modules struct {
	K8sClient  *K8sClient  `json:"k8s_client"`
	PromServer *PromServer `json:"prom_server"`
}

type CloudAgentConfig struct {
	// cloudagentagent中的模块配置
	Modules *Modules `json:"modules,omitempty"`
}

type K8sClient struct {
	Enable             bool               `json:"enable"`
	MasterLBIp         string             `json:"master_lb_ip"`
	MasterLBPort       int                `json:"master_lb_port"`
	RedisIp            string             `json:"redis_ip"`
	RedisPort          int                `json:"redis_port"`
	PodInfo            *v1.Pod            `json:"pod_info"`
	DeploymentInfo     *appsv1.Deployment `json:"deployment_info"`
	KubeConfigFilePath string             `json:"kube_config_file_path"`
	DecoderBufferSize  int                `json:"decoder_buffer_size"`
}

type PromServer struct {
	Enable bool `json:"enable"`
}
