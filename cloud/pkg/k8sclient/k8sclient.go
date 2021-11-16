package k8sclient

import (
	"github.com/gomodule/redigo/redis"
	"github.com/wonderivan/logger"
	v1 "k8s.io/api/core/v1"
	"keep/cloud/pkg/common"
	k8sclientconfig "keep/cloud/pkg/k8sclient/config"
	"keep/core"
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"os"
	"strconv"
	"time"
)

type K8sClient struct {
	enable       bool
	MasterLBIp   string
	MasterLBPort int
	PodInfo      *v1.Pod
}

// Register 注册healthzagent模块
func Register(k *cloudagent.K8sClient) {
	k8sclientconfig.InitConfigure(k)
	k8sclient, err := NewK8sClient(k.Enable)
	if err != nil {
		logger.Error("初始化k8sclient失败...:", err)
		os.Exit(1)
		return
	}
	core.Register(k8sclient)
}

func (k *K8sClient) Name() string {
	return modules.K8sClientModule
}

func (k *K8sClient) Group() string {
	return modules.K8sClientGroup
}

//Enable indicates whether this module is enabled
func (k *K8sClient) Enable() bool {
	return k.enable
}

func (k *K8sClient) Start() {
	logger.Debug("k8sclient开始启动....")
	// 检查k8s集群apiserver状态

	// 检查redis在线状态 如果不在线就由naive_engine 在master集群中创建pod
	checkRedisAliveness()
	for {
		time.Sleep(time.Second)
	}
}

func NewK8sClient(enable bool) (*K8sClient, error) {
	return &K8sClient{
		enable: enable,
	}, nil
}

func checkRedisAliveness() {
	port := k8sclientconfig.Config.RedisPort
	_, err := redis.Dial("tcp", k8sclientconfig.Config.RedisIp+":"+strconv.Itoa(port))
	if err != nil {
		logger.Fatal("redis初始化失败...,err:", err)
		os.Exit(1)
	}
	logger.Debug("redis初始化成功...")
}
