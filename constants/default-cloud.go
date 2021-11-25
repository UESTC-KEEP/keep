package constants

const CloudAgentName = "CloudAgent"
const CloudConfigeFilesSourceDir = "../../../cloud/shells/confs/keepedge"

// 全局配置
const (
	DefaultKubeConfigPath  = "/root/.kube/config"
	DefaultCloudConfigFile = "/etc/keepedge/config/cloudagent.yml"

	NodeName = "NodeName"
)

// k8sclient配置
const (
	// DefaultMasterLBIp master集群的负载均衡ip 若是单master集群就是用masterip即可
	DefaultMasterLBIp       = "192.168.1.140"
	DefaultMasterLBPort     = 6443
	DefaultRedisServerIp    = "192.168.1.140"
	DefaultRedisServerPort  = 32379
	DefaultRedisConfigMap   = "/etc/keepedge/ymls/redis-standalone-conf.yml"
	DefaultRedisSVC         = "/etc/keepedge/ymls/redis-svc.yml"
	DefaultRedisStatefulSet = "/etc/keepedge/ymls/redis-statefulset.yml"
	DefaultNameSpace        = "default"
)

const (
	DefaultCAURL   = "/ca.crt"
	DefaultCertURL = "/edge.crt"
)
