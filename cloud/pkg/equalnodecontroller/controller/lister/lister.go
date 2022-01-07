package lister

import (
	"context"
	"github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/equalnode/v1alpha1"
	"github.com/UESTC-KEEP/keep/cloud/pkg/equalnodecontroller/config"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetAllEqnd() *v1alpha1.EqualNodeList {
	eqndlist, err := config.EqndClient.KeepedgeV1alpha1().EqualNodes("default").List(context.Background(), metav1.ListOptions{})
	//eqndlist, err := config.EqndClient.KeepedgeV1().EqualNodes(corev1.NamespaceAll).List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		logger.Error(err)
		return nil
	}
	return eqndlist
}
