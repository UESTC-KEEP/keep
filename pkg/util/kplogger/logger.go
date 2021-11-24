package kplogger

import (
	"encoding/json"
	"fmt"
	"keep/constants"
	"os"
	"path/filepath"

	"github.com/wonderivan/logger"
)

const LOG_TAG = "<LOGGER>:"

func loadDefaultLoggerConf() *logConfig { //读default
	p_cfg := &logConfig{ //TODO 以后写在default里面
		Console: &consoleLogger{
			Level:    "TRAC",
			Colorful: true,
		},
		File: &fileLogger{
			Filename:   constants.DefaultEdgeLogFiles,
			Level:      "TRAC",
			Daily:      true,
			MaxLines:   1000000,
			MaxSize:    1,
			MaxDays:    -1,
			Append:     true,
			PermitMask: "0600",
		},
		Conn: &connLogger{},
	}
	return p_cfg
}

//把配置写入文件中，供logger包读取
func WirteLoggerConfigToFile(file_path string, p_cfg *logConfig) { //覆盖
	data, err := json.Marshal(p_cfg)

	if nil != err {
		Error(LOG_TAG + "cannont   Marshal config :" + err.Error())
	} else {
		Infof("%s, data=%s", LOG_TAG, string(data))
		err := os.WriteFile(file_path, data, constants.DefaultUnixFilePermit)
		if nil != err {
			Error(LOG_TAG + "cannont  write logger config file:" + err.Error())
		}
	}
}

func CreateDefaultLoggerConfigFile(file_path string) {
	p_default_cfg := loadDefaultLoggerConf()
	WirteLoggerConfigToFile(file_path, p_default_cfg)
}

func CreateLogFile(file_path string) {
	fp, err := os.Create(file_path)
	if nil == err {
		fp.Close()
	} else {
		Error(LOG_TAG + "Failed to   create log file :" + err.Error())
	}
}

func CheckAndCreateFile(file_path string, createor_handle func(string)) { //TODO 也许应该放utils里面
	if (0 == len(file_path)) || ('/' != file_path[0]) { //不支持windows的路径格式 TODO
		return
	}

	_, err := os.Stat(file_path)
	if nil == err {
		return
	}

	Warnf("%s:cannont find %s,creating a new one now", LOG_TAG, file_path)

	paths, _ := filepath.Split(file_path)
	err = os.MkdirAll(paths, constants.DefaultUnixDirectoryPermit)
	if nil == err {
		createor_handle(file_path)
	} else {
		Errorf("%s: Cannont create directory \"%s\",err=%s", LOG_TAG, paths, err.Error())
	}
}

func InitKeepLogger() {

	CheckAndCreateFile(constants.DefaultEdgeLogFiles, CreateLogFile)

	logger_cfg_path := constants.DefaultEdgeLoggerConfFile
	CheckAndCreateFile(logger_cfg_path, CreateDefaultLoggerConfigFile)
	err := logger.SetLogger(logger_cfg_path)
	if err != nil {
		Error(LOG_TAG + "Keeploger初始化失败:" + err.Error())
		return
	}
	Debug(LOG_TAG + "Keeploger初始化成功...")

}

func Painc(f interface{}, v ...interface{}) {
	logger.Painc(f, v...)
}

func Paincf(format string, a ...interface{}) {
	logger.Painc(fmt.Sprintf(format, a...))
}

func Fatal(f interface{}, v ...interface{}) {
	logger.Fatal(f, v...)
}

func Fatalf(format string, a ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, a...))
}

func Emer(f interface{}, v ...interface{}) {
	logger.Emer(f, v...)
}

func Emerf(format string, a ...interface{}) {
	logger.Emer(fmt.Sprintf(format, a...))
}

func Alert(f interface{}, v ...interface{}) {
	logger.Alert(f, v...)
}

func Alertf(format string, a ...interface{}) {
	logger.Alert(fmt.Sprintf(format, a...))
}

func Crit(f interface{}, v ...interface{}) {
	logger.Crit(f, v...)
}

func Critf(format string, a ...interface{}) {
	logger.Crit(fmt.Sprintf(format, a...))
}

func Error(f interface{}, v ...interface{}) {
	logger.Error(f, v...)
}

func Errorf(format string, a ...interface{}) {
	logger.Error(fmt.Sprintf(format, a...))
}

func Warn(f interface{}, v ...interface{}) {
	logger.Warn(f, v...)
}

func Warnf(format string, a ...interface{}) {
	logger.Warn(fmt.Sprintf(format, a...))
}

func Info(f interface{}, v ...interface{}) {
	logger.Info(f, v...)
}

func Infof(format string, a ...interface{}) {
	logger.Info(fmt.Sprintf(format, a...))
}

func Debug(f interface{}, v ...interface{}) {
	logger.Debug(f, v...)
}

func Debugf(format string, a ...interface{}) {
	logger.Debug(fmt.Sprintf(format, a...))
}

func Trace(f interface{}, v ...interface{}) {
	logger.Trace(f, v...)
}

func Tracef(format string, a ...interface{}) {
	logger.Trace(fmt.Sprintf(format, a...))
}
