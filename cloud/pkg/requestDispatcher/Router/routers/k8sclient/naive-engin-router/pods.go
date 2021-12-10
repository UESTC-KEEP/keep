package naive_engin_router

import (
	"keep/pkg/util/core/model"
)

type Pods struct {
	Operation Operation
}

type Operation struct {
	List model.MessageRoute
}
