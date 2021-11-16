// Package tests 用于测试本模块中的所有代码
package tests

import (
	"fmt"
	"keep/edge/pkg/healthzagent/server"
	"testing"
)

func TestGetCpuStatus(t *testing.T) {
	fmt.Println(server.GetCpuStatus())
}

func TestGetMemStatus(t *testing.T) {
	info := server.GetMemStatus()
	fmt.Println("内存占比：", float64(info.Used)/float64(info.Total))
}

func TestGetBasicStatus(t *testing.T) {
	fmt.Println(server.GetBasicStatus())
}

func TestGetDiskStatus(t *testing.T) {
	server.GetDiskStatus()
}

func TestGetNetIOStatus(t *testing.T) {
	server.GetNetIOStatus()
}
