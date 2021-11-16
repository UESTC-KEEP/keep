package test

import (
	"github.com/wonderivan/logger"
	"keep/conf"
	"testing"
)

func TestGetConfig(t *testing.T) {
	logger.Debug(conf.GetStringConfig("dbhost"))
	//fmt.Println(GetStringConfig("dbhost"))
}
