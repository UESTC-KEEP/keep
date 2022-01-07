package controller

import (
	crdv1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/equalnode/v1alpha1"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"reflect"
)

func (eqndctl *EqualNodeController) equalNodeUpated(eqnd *crdv1.EqualNode) {
	value, ok := eqndctl.equalnodeManager.EqualNode.Load(eqnd.Name)
	eqndctl.equalnodeManager.EqualNode.Store(eqnd.Name, eqnd)
	if ok {
		cacheEqualNode := value.(*crdv1.EqualNode)
		if isEqualNodeUpdated(cacheEqualNode, eqnd) {
			logger.Info("----------- crd更新:  ")
			logger.Info("----------- crd更新:  ", eqnd)
		}
	}
}

// isDeviceUpdated 检查eqnd crd是否更新
func isEqualNodeUpdated(oldeqnd *crdv1.EqualNode, neweqnd *crdv1.EqualNode) bool {
	// does not care fields
	oldeqnd.ObjectMeta.ResourceVersion = neweqnd.ObjectMeta.ResourceVersion
	oldeqnd.ObjectMeta.Generation = neweqnd.ObjectMeta.Generation
	// return true if ObjectMeta or Spec or Status changed, else false
	return !reflect.DeepEqual(oldeqnd.ObjectMeta, neweqnd.ObjectMeta) || !reflect.DeepEqual(oldeqnd.Spec, neweqnd.Spec)
}
