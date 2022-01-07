package k8sclientrouter

import (
	kubedge_engin_router "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient/kubedge-engin-router"
	naiveenginrouter "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient/naive-engin-router"
)

type K8sClientRouter struct {
	NaiveEngine    naiveenginrouter.NaiveEngine
	KubeedgeEngine kubedge_engin_router.KubeedgeEngine
}
