package k8sclientrouter

import (
	naiveenginrouter "keep/cloud/pkg/requestDispatcher/Router/routers/k8sclient/naive-engin-router"
)

type K8sClientRouter struct {
	NaiveEngine naiveenginrouter.NaiveEngine
}
