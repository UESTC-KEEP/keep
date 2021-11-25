package crd_engin

import (
	naive_engine "keep/cloud/pkg/k8sclient/naive-engine"
	"keep/constants"
	"keep/pkg/util/kplogger"
)

type CrdEngineImpl struct{}

func NewCrdEngineImpl() *CrdEngineImpl {
	return &CrdEngineImpl{}
}


func (cei *CrdEngineImpl) CreateCrd(Dir string) error {
	var err error
	if err != nil {
		kplogger.Error(err)
	}
	naive_engine.CreatResourcesByYAML(constants.DefaultCrdsDir+"/"+"stuCrd.yaml", constants.DefaultNameSpace)

	naive_engine.CreatResourcesByYAML(constants.DefaultCrdsDir+"/"+"new-student.yaml", constants.DefaultNameSpace)
	//Clientset.Students("default").Delete(context.TODO(),"new-student",metav1.DeleteOptions{})
	return err
}
