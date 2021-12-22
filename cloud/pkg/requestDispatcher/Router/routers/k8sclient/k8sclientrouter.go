package k8sclientrouter

import (
	kubedge_engin_router "keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient/kubedge-engin-router"
	naiveenginrouter "keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient/naive-engin-router"
)

type K8sClientRouter struct {
	NaiveEngine    naiveenginrouter.NaiveEngine
	KubeedgeEngine kubedge_engin_router.KubeedgeEngine
}
