package cloud

import "keep/constants"

// NewDefaultEdgeAgentConfig returns a full EdgeCoreConfig object
func NewDefaultEdgeAgentConfig() *CloudAgentConfig {
	return &CloudAgentConfig{
		Modules: &Modules{
			K8sClient: &K8sClient{
				Enable:       true,
				MasterLBIp:   constants.DefaultMasterLBIp,
				MasterLBPort: constants.DefaultMasterLBPort,
				RedisIp:      constants.DefaultRedisServerIp,
				RedisPort:    constants.DefaultRedisServerPort,
				PodInfo:      nil,
			},
		},
	}
}
