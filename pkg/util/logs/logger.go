package logs

import (
	"github.com/wonderivan/logger"
)

func InitKeepLogger() {
	//logger.SetLogger(`{"Console": {"level": "TRAC"}}`)
	err := logger.SetLogger("/etc/keepedge/logger_conf.json")
	if err != nil {
		logger.Error("Keeploger初始化失败:" + err.Error())
		return
	}
	logger.Debug("Keeploger初始化成功...")
	return
}
