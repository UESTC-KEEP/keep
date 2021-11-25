package kplogger

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	TestLoadDefaultLoggerConf(t)
}

func TestLoadDefaultLoggerConf(t *testing.T) {
	cfg := loadDefaultLoggerConf()
	data, err := json.Marshal(cfg)
	if nil != err {
		fmt.Println("log test err:", err)
	} else {
		fmt.Println(string(data))
	}
}
