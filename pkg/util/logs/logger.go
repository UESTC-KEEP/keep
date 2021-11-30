package logs

import (
	"keep/pkg/util/loggerv1.0.1"
)

func InitKeepLogger() {
	//logger.SetLogger(`{"Console": {"level": "TRAC"}}`)
	err := logger.SetLogger("../../../pkg/util/logs/logger_conf.json")
	if err != nil {
		logger.Error("Keeploger初始化失败:" + err.Error())
		return
	}
	logger.Debug("Keeploger初始化成功...")
	return
}
