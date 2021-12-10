package routers

import (
	k8sclientrouter "keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient"
	naive_engin_router "keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient/naive-engin-router"
	"keep/pkg/util/core/model"
	"sync"
)

// "$uestc/keep/k8sclient/naiveengine/pods/"

type KeepRouters struct {
	K8sClientRouter k8sclientrouter.K8sClientRouter
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
						Operation: naive_engin_router.Operation{
							List: model.MessageRoute{
								Operation: "list",
								Resource:  "$uestc/keep/k8sclient/naiveengine/pods/",
							},
						},
					},
				},
			},
		}
	})
}
