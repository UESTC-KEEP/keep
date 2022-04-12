package routers

import (
	k8sclientrouter "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient"
	kubedge_engin_router "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient/kubedge-engin-router"
	naive_engin_router "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient/naive-engin-router"
	kafka_router "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router/routers/kafka"
	kafka_metrics "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router/routers/kafka/metrics"
	"sync"
)

// "$uestc/keep/k8sclient/naiveengine/pods/"

type KeepRouters struct {
	K8sClientRouter k8sclientrouter.K8sClientRouter
	Kafka           kafka_router.FluentKafkaRouter
}

var KeepRouter KeepRouters
var once sync.Once

func InitRouters() {
	// 初始化router路由函数绑定
	once.Do(func() {
		KeepRouter = KeepRouters{
			K8sClientRouter: k8sclientrouter.K8sClientRouter{
				NaiveEngine: naive_engin_router.NaiveEngine{
					Pods: naive_engin_router.Pods{
						Resources: "$uestc/keep/k8sclient/naiveengine/pods/",
						Operation: naive_engin_router.Operation{
							List: "list",
						},
					},
				},
				KubeedgeEngine: kubedge_engin_router.KubeedgeEngine{
					Devices: kubedge_engin_router.Devices{
						Resources: "$uestc/keep/k8sclient/kubeedgeengin/devices/",
						Operation: kubedge_engin_router.Operation{
							List: "list",
						},
					},
				},
			},
			Kafka: kafka_router.FluentKafkaRouter{
				Metrics: kafka_metrics.Metrics{
					Resources: "uestc-keep-kafka-metrics",
					Operation: kafka_metrics.Operation{
						Push: "push",
					},
				},
			},
		}
	})
}
