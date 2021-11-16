package constants

// 全局配置
const (
	DefaultKubeConfigPath  = "/root/.kube/config"
	DefaultCloudConfigFile = "/etc/keepedge/config/cloudagent.yml"
)

// k8sclient配置
const (
	// DefaultMasterLBIp master集群的负载均衡ip 若是单master集群就是用masterip即可
	DefaultMasterLBIp   = "192.168.1.140"
	DefaultMasterLBPort = 6443
)
