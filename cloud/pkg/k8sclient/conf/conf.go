package conf

import (
	cloudagent "keep/pkg/apis/compoenentconfig/cloudagent/v1alpha1"
	"sync"


)


var Config K8sConfigure
var once sync.Once
type K8sConfigure struct {
	cloudagent.K8sClient
}
func InitConfig(kc *cloudagent.K8sClient){
	once.Do(func() {
		Config= K8sConfigure{
			K8sClient:*kc,
		}
	},
	)
}
func Get() *K8sConfigure{
	return &Config
}

