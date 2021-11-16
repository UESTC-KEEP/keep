package conf

import (
	"github.com/astaxie/beego/config"
	"github.com/wonderivan/logger"
)

func GetStringConfig(configName string) string {
	conf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		logger.Error(err.Error())
		//fmt.Println(err.Error())
	}
	return conf.String(configName)
}

func GetIntConfig(configName string) int {
	conf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		logger.Error(err.Error())
		//fmt.Println(err.Error())
	}
	res, _ := conf.Int(configName)
	if err != nil {
		return -1
	}
	return res
}
