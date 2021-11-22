package constants

const EdgeAgentName = "EdgeAgent"
const EdgeConfigeFilesSourceDir = "../../../edge/shells/confs/keepedge"

const (
	KeepBasepath     = "/etc/keepedge/"
	KeepBaseConfPath = KeepBasepath + "config/"
	KeepBaseLogPath  = "/var/log/keepedge/"
)

const (
	// DefaultConfigDir 默认的存储文件位置
	DefaultEdgeagentConfigFile = KeepBaseConfPath + "/edgecore.yml"
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
	DefaultHttpServer    = "http://192.168.1.140"
	DefaultCloudHttpPort = 20000
	DefaultEdgeHeartBeat = 15
	DefaultEdgePort      = 20350
	DefaultLogsTopic     = "keep_log_topic"
	DefaultDataTopic     = "keep_data_topic"
	// DefaultLogsQueenSize 日志缓冲队列默认大小
	DefaultLogsQueenSize = 100
	// DefaultDataQueenSize 数据缓冲队列默认大小
	DefaultDataQueenSize = 100
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
