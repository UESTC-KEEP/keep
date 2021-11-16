package logs

import (
	"github.com/wonderivan/logger"
)

func InitKeepLogger() {
	err := logger.SetLogger("/etc/keepedge/logger_conf.json")
	if err != nil {
		logger.Error("cloudloger初始化失败:" + err.Error())
		return
	}
	logger.Debug("cloudloger初始化成功...")
	return
}
