package edge

const EdgeAgentName = "EdgeAgent"
const EdgeConfigeFilesSourceDir = "../../../edge/shells/confs/keepedge"

const (
	KeepBasepath              = "/etc/keepedge/"
	DefaultEdgecoreConfigFile = "/etc/kubeedge/config/edgecore.yaml"
	KeepBaseConfPath          = KeepBasepath + "config/"
	KeepBaseLogPath           = "/var/log/keepedge/"
	// DefaultBeehiveTimeout	等待接收消息时间上限 ms
	DefaultBeehiveTimeout = 500
)

const (
	DefaultUnixDirectoryPermit = 0660 //不可执行，非本组成员不能访问
	DefaultUnixFilePermit      = 0660 //不可执行，非本组成员不能访问
)

const (
	// DefaultEdgeagentConfigFile  默认的存储文件位置
	DefaultEdgeagentConfigFile = EdgeConfigeFilesSourceDir + "/config/edgeagent.yaml"
	//KeepEdgeVersion	keepedge 版本信息
	KeepEdgeVersion = "0.0.1"
)

// healthzagent 全局静态配置
const (
	DefaultEdgeHealthInterval = 30
)

//
const (
	DefaultEdgeLogFiles       = KeepBaseLogPath + "keep_edgeagent_logs.log"
	DefaultCloudLogFiles      = KeepBaseLogPath + "keep_edgeagent_logs.log"
	DefaultEdgeLoggerConfFile = KeepBaseConfPath + "logger_conf.json"
)

// EdgePublisher 全局配置
const (
	DefaultHttpServer    = "192.168.1.128"
	DefaultCloudHttpPort = 20001
	DefaultWebSocketPort = 20000
	DefaultEdgeHeartBeat = 15
	DefaultEdgePort      = 20350
	DefaultLogsTopic     = "keep_log_topic"
	DefaultDataTopic     = "keep_data_topic"
	// DefaultLogsQueenSize 日志缓冲队列默认大小
	DefaultLogsQueenSize = 100
	// DefaultDataQueenSize 数据缓冲队列默认大小
	DefaultDataQueenSize = 100
	// DefaultMetricsPort metrics暴露端口
	DefaultMetricsPort       = 8080
	DefaultMsgSendBufferSize = 500
)

// mqtt配置
const (
	DefaultTestingMQTTServer = "192.168.1.40"
	DefaultTestingMQTTPort   = 1883
	DefaultMqttChanSize      = 4
	DefaultMqttLogTag        = "<MQTT>"
	// DefaultDeviceMqttTopics 默认想要监听的设备mqtt主题 以; 进行分割多个主题
	DefaultDeviceMqttTopics = "clock_sensor"
)

// EdgeTwin配置
const (
	DefaultEdgeTwinSqliteFilePath = "/var/lib/keepedge/edgeagent.db"
)

const (
	// Certificates
	DefaultConfigDir = "/etc/keepedge/config/"
	DefaultCAFile    = "/etc/keepedge/ca/rootCA.crt"
	DefaultCAKeyFile = "/etc/keepedge/ca/rootCA.key"
	DefaultCertFile  = "/etc/keepedge/certs/server.crt"
	DefaultKeyFile   = "/etc/keepedge/certs/server.key"
)
