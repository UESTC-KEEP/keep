package controller

import (
	"context"
	crdv1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/equalnode/v1alpha1"
	"github.com/UESTC-KEEP/keep/cloud/pkg/common/client"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (eqndctl *EqualNodeController) equalNodeDeleted(eqnd *crdv1.EqualNode) {
	eqndctl.equalnodeManager.EqualNode.Delete(eqnd.Name)
	logger.Info("----------- crd删除:  ")
	logger.Info("----------- crd删除:  ", eqnd)

	err := client.GetKubeClient().AppsV1().Deployments(apiv1.NamespaceDefault).Delete(context.Background(), "keep-eqnd-test-nginx", metav1.DeleteOptions{})
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Warn("删除Deployments " + "keep-eqnd-test-nginx " + "成功....")
}
