package config

import "sync"

var Config Configure
var once sync.Once

type Configure struct {
	// v1alpha1.RequestDispatcher
	//KubeAPIConfig *v1alpha1.KubeAPIConfig
	Ca    []byte
	CaKey []byte
	Cert  []byte
	Key   []byte
}
