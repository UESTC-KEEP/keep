// Package tests 用于测试本模块中的所有代码
package tests

import (
	"fmt"
	"github.com/UESTC-KEEP/keep/edge/pkg/healthzagent/server"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"io/ioutil"
	"net/http"
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

func TestGetDiskStorageStatus(t *testing.T) {
	info := server.GetDiskStorageStatus()
	for name, obj := range *info {
		fmt.Println(name, "\t", obj)
	}

}

func TestGetDiskIOStatus(t *testing.T) {
	info := server.GetDiskIOStatus()
	for name, obj := range *info {
		fmt.Println(name, "\t", obj)
	}
}

func TestGetNetIOStatus(t *testing.T) {
	server.GetNetIOStatus()
}

func TestNodeExporter(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:9100/metrics")
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
