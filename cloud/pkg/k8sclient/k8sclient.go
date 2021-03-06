package k8sclient

import (
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/modules"
	k8sclientconfig "github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/config"
	naive_engine "github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/naive-engine"
	"github.com/UESTC-KEEP/keep/cloud/pkg/k8sclient/watchengine"
	"github.com/UESTC-KEEP/keep/constants/cloud"
	cloudagent "github.com/UESTC-KEEP/keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"github.com/UESTC-KEEP/keep/pkg/util/core"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"github.com/gomodule/redigo/redis"
	v1 "k8s.io/api/core/v1"
	"os"
	"strconv"
	"strings"
	"sync"
)

type K8sClient struct {
	enable       bool
	MasterLBIp   string
	MasterLBPort int
	PodInfo      *v1.Pod
}

// Register 注册K8sClient模块
func Register(k *cloudagent.K8sClient) {
	k8sclientconfig.InitConfigure(k)
	k8sclient, err := NewK8sClient(k.Enable)
	if err != nil {
		logger.Error("初始化k8sclient失败...:", err)
		os.Exit(1)
	}
	core.Register(k8sclient)
}

func (k *K8sClient) Cleanup() {
	//logger.Debug("准备清理模块：",modules.K8sClientModule)
	deleteRedis()
	//naive_engine.DeleteResourceByYAML(constants.DefaultCrdsDir+"/"+"stuCrd.yaml", constants.DefaultNameSpace)
	//naive_engine.DeleteResourceByYAML(constants.DefaultCrdsDir+"/"+"new-student.yaml", constants.DefaultNameSpace)
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
	// 初始化keepedge cloud环境
	initKeepEdgeEnv()
	//metrics_server.NewMetricServerImpl().CheckCadvisorStatus([]string{"192.168.1.140:6443", "192.168.1.141:6443"})
	//var count int
	//for {
	//	count++
	//	fmt.Println(count)
	//	time.Sleep(time.Second)
	//}
	// 查询所有的device
	//go kubeedge_engine.NewKubeEdgeEngine().GetDevicesByNodeName("")
	// 启动系统需要的所有informers们
	go watchengine.StartAllInformers()
	//	所有项目准备完成启动路由
	go StartK8sClientRouter()
}

func NewK8sClient(enable bool) (*K8sClient, error) {
	return &K8sClient{
		enable: enable,
	}, nil
}

func initKeepEdgeEnv() {
	// 检查k8s版本是否支持
	go checkK8sVersion()
	// 检查有没有keepedge的namesopace
	// 没有就创建
	//go checkNamespaceAliveness()
	// 检查k8s集群apiserver状态
	// 检查redis在线状态 如果不在线就由naive_engine 在master集群中创建statefulset
	//go checkRedisAliveness()
	// 创建crd
	//go checkeepCrd()

}

// 检查constants.DefaultKeepEdgeNameSpace是否存在不存在就创建之
func checkNamespaceAliveness() {
	ns, _ := naive_engine.NewNaiveEngine().GetNamespaceByName(cloud.DefaultKeepEdgeNameSpace)
	if ns == nil {
		// 不存在需要的ns就创建之
		ns_, err := naive_engine.NewNaiveEngine().CreateNamespaceByName(cloud.DefaultKeepEdgeNameSpace)
		if err != nil {
			logger.Fatal("创建ns失败：", err)
		}
		logger.Debug("创建的ns：", ns_)
		return
	}
	logger.Debug("查询到ns: ", ns)
}

// 检查集群中redis状态  无法通信就拉起redis
func checkRedisAliveness() {
	port := k8sclientconfig.Config.RedisPort
	_, err := redis.Dial("tcp", k8sclientconfig.Config.RedisIp+":"+strconv.Itoa(port))
	if err != nil {
		logger.Warn("redis初始化失败...,err:", err)
		logger.Debug("准备拉起redis...")
		// 创建redis的configMap
		naive_engine_impl := naive_engine.NewNaiveEngine()
		naive_engine_impl.CreatResourcesByYAML(cloud.DefaultRedisConfigMap, cloud.DefaultNameSpace)
		// 创建redis服务
		naive_engine_impl.CreatResourcesByYAML(cloud.DefaultRedisSVC, cloud.DefaultNameSpace)
		// 创建redis statfulset
		naive_engine_impl.CreatResourcesByYAML(cloud.DefaultRedisStatefulSet, cloud.DefaultNameSpace)
	} else {
		logger.Debug("redis 在线运行中.....")
	}
	logger.Debug("redis初始化成功...")
}

// deleteRedis 删除redis各组件
func deleteRedis() {
	ns := cloud.DefaultNameSpace
	var compoenents = []string{cloud.DefaultRedisConfigMap, cloud.DefaultRedisSVC, cloud.DefaultRedisStatefulSet}
	var wg sync.WaitGroup
	for i := 0; i < len(compoenents); i++ {
		wg.Add(1)
		go func(i int) {
			logger.Debug("开始删除redis组件 " + compoenents[i] + " ...")
			err := naive_engine.NewNaiveEngine().DeleteResourceByYAML(compoenents[i], ns)
			if err != nil {
				wg.Done()
				logger.Error(err)
				return
			} else {
				wg.Done()
			}
		}(i)
	}
	wg.Wait()
}

func checkK8sVersion() {
	versionstr := naive_engine.NewNaiveEngine().GetK8sVersion()
	if versionstr == "" {
		logger.Error("获取kubernetes版本失败...")
		return
	}
	version := strings.Split(versionstr, ".")

	// 只检查大版本不查小版本
	for _, surpport := range strings.Split(cloud.DefaultKubeEdgeSupportedK8sVersion, ";") {
		if version[1] == surpport {
			logger.Debug("kubernetes版本：" + naive_engine.NewNaiveEngine().GetK8sVersion())
			return
		}
	}
	beehiveContext.Cancel()
	logger.Fatal("kubernetes version:", naive_engine.NewNaiveEngine().GetK8sVersion(), " 不受kubeedge支持无法启动...")
}

// 拉起keep crd
func checkeepCrd() {
	// 创建crd
	naive_engine.NewNaiveEngine().CreatResourcesByYAML(cloud.DefaultKeepEqndCrd, cloud.DefualtKeepNamespace)
	naive_engine.NewNaiveEngine().CreatResourcesByYAML(cloud.DefaultKeepTrqCrd, cloud.DefualtKeepNamespace)
	naive_engine.NewNaiveEngine().CreatResourcesByYAML(cloud.DefualtKeepTenantCrd, cloud.DefualtKeepNamespace)
}

// 拉起service
