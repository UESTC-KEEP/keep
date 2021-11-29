package constants

const CloudAgentName = "CloudAgent"
const CloudConfigeFilesSourceDir = "../../../cloud/shells/confs/keepedge"

// 全局配置
const (
	DefaultKubeConfigPath    = "/home/et/.kube/config"
	DefaultCloudConfigFile   = "/etc/keepedge/config/cloudagent.yml"
	DefaultDecoderBufferSize = 100
	NodeName                 = "NodeName"
)

// k8sclient配置
const (
	// DefaultMasterLBIp master集群的负载均衡ip 若是单master集群就是用masterip即可
	DefaultMasterLBIp = "192.168.1.140"
	// DefaultMasterIpPort 默认需要进行监控的主机ip:port 多个master主机使用;间隔
	DefaultMasterIpPort = DefaultMasterLBIp + ":6443"
	// DefaultMasterMetricTimeout 获取每台master  metrics信息的超时时间
	DefaultMasterMetricTimeout = 5000
	DefaultMasterLBPort        = 6443
	DefaultRedisServerIp       = DefaultMasterLBIp
	DefaultRedisServerPort     = 32379
	DefaultRedisConfigMap      = "/etc/keepedge/ymls/redis-standalone-conf.yaml"
	DefaultRedisSVC            = "/etc/keepedge/ymls/redis-svc.yaml"
	DefaultRedisStatefulSet    = "/etc/keepedge/ymls/redis-statefulset.yaml"
	DefaultNameSpace           = "default"
	DefaultCrdsDir             = "/etc/keepedge/keep-crds"
)

// requestDispatcher配置
const (
	DefaultHTTPPort      = 20001
	DefaultWebSocketPort = 20000
	DefaultCAURL         = "/ca.crt"
	DefaultCertURL       = "/edge.crt"
	DefaultWebSocketUrl  = "/v1/keepedge/connect"

	SessionKeyHostNameOverride = "SessionHostNameOverride"
	SessionKeyInternalIP       = "SessionInternalIP"
)

const (
	DefaultPromServerMetricsPort = 20080
)

// equalnodeController配置
const (
	DefaultMasterURL       = DefaultMasterLBIp
	DefaultKubeConfig      = DefaultKubeConfigPath
	DefaultAlsoLogToStdErr = true
)

//LogPublisher
const Url = DefaultMasterLBIp + ":4560"
const ContentType = "apllication/json;charset=utf-8"
