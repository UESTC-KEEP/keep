package config

import (
	cloudagent "keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
	"sync"
)

var Config Configure
var once sync.Once

type Configure struct {
	cloudagent.RequestDispatcher
	//KubeAPIConfig *v1alpha1.KubeAPIConfig
	Ca    []byte // ca证书 包含公私钥
	CaKey []byte // 申请CA证书的 私钥
	Cert  []byte // server CA证书
	Key   []byte // server 私钥
	Token string
}

func InitConfigure(r *cloudagent.RequestDispatcher) {
	once.Do(func() {
		Config = Configure{
			RequestDispatcher: *r,
		}
	})
}
