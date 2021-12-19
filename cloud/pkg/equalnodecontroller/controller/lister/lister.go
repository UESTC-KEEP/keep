package lister

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "keep/cloud/pkg/apis/keepedge/v1"
	"keep/cloud/pkg/equalnodecontroller/config"
	logger "keep/pkg/util/loggerv1.0.1"
)

func GetAllEqnd() *v1.EqualNodeList {
	eqndlist, err := config.EqndClient.KeepedgeV1().EqualNodes("default").List(context.Background(), metav1.ListOptions{})
	//config.EqndClient.KeepedgeV1().EqualNodes(corev1.NamespaceAll).List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		logger.Error(err)
		return nil
	}
	return eqndlist
}
