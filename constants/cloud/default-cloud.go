package cloud

const CloudAgentName = "CloudAgent"
const CloudConfigeFilesSourceDir = "../../../cloud/shells/confs/keepedge"

// 全局配置
const (
	DefaultKeepEdgeNameSpace = "keepedge"
	DefaultKeepConfigDir     = "/etc/keepedge/"
	DefaultKubeConfigPath    = "/home/et/.kube/config"
	DefaultCloudConfigFile   = DefaultKeepConfigDir + "config/cloudagent.yml"
	DefaultKeepCrd           = DefaultKeepConfigDir + "ymls/keepcrd/equivalentnode.yaml"
	DefualtKeepNamespace     = "keepedge"
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
	DefaultRedisConfigMap      = DefaultKeepConfigDir + "ymls/redis-standalone-conf.yaml"
	DefaultRedisSVC            = DefaultKeepConfigDir + "ymls/redis-svc.yaml"
	DefaultRedisStatefulSet    = DefaultKeepConfigDir + "ymls/redis-statefulset.yaml"
	DefaultNameSpace           = "default"
)

// requestDispatcher配置
const (
	DefaultHTTPPort             = 20001
	DefaultWebSocketPort        = 20000
	DefaultCAURL                = "/ca.crt"
	DefaultCertURL              = "/edge.crt"
	DefaultWebSocketUrl         = "/v1/keepedge/connect"
	DefaultKeepCloudIP          = "192.168.1.121"
	DefaultTokenRefreshDuration = 12
	SessionKeyHostNameOverride  = "SessionHostNameOverride"
	SessionKeyInternalIP        = "SessionInternalIP"
)

const (
	DefaultPromServerMetricsPort = 20080
)

// equalnodeController配置
const (
	DefaultMasterURL       = DefaultMasterLBIp
	DefaultAlsoLogToStdErr = true
)

//LogPublisher
const Url = DefaultMasterLBIp + ":4560"
const ContentType = "apllication/json;charset=utf-8"

var Address = []string{"192.168.1.103:9092", "192.168.1.103:9093"}

const (
	OrginTopic = "topic"
	ParseTopic = "topic1"
)
