package constants

const (
	// 默认的存储文件位置
	DefaultConfigDir = "/etc/keepedge/config/"
)

// healthzagent 全局静态配置
const (
	DefaultEdgeHealthInterval     = 10
	DefaultHealthzToCloudInterval = 60
	DefaultMqttCacheQueueSize     = 10
)

//
const (
	DefaultLogFiles = "/var/log/keepedge/keep_edgeagent_logs.log$$$$/var/log/test.log"
)

// EdgePublisher 全局配置
const (
	DefaultHttpServer    = "http://192.168.1.140"
	DefaultHttpPort      = 20000
	DefaultEdgeHeartBeat = 15
	DefaultServePort     = 20350
	DefaultLogsTopic     = "keep_log_topic"
	DefaultDataTopic     = "keep_data_topic"
	// 日志缓冲队列默认大小
	DefaultLogsQueenSize = 100
	// 数据缓冲队列默认大小
	DefaultDataQueenSize = 100
)
