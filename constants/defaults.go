package constants

const (
	// DefaultConfigDir 默认的存储文件位置
	DefaultEdgeagentConfigFile = "/etc/keepedge/config/edgecore.yml"
	//KeepEdgeVersion	keepedge 版本信息
	KeepEdgeVersion = "0.0.1"
)

// healthzagent 全局静态配置
const (
	DefaultEdgeHealthInterval     = 30
	DefaultHealthzToCloudInterval = 60
	DefaultMqttCacheQueueSize     = 10
)

//
const (
<<<<<<< HEAD
	DefaultLogFiles = "/var/logs/keepedge/keep_edgeagent_logs.logs$$$$/var/logs/test.logs"
=======
	DefaultEdgeLogFiles  = "/var/log/keepedge/keep_edgeagent_logs.log"
	DefaultCloudLogFiles = "/var/log/keepedge/keep_edgeagent_logs.log"
>>>>>>> bddbd7e0f200a771b61cbb6932118d2c7492d2c4
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
