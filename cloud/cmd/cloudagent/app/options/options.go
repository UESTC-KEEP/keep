package options

import (
	"keep/constants"
	"path"
)

func NewDefaultCloudAgentOptions() *EdgeAgentOptions {
	return &EdgeAgentOptions{
		ConfigFile: path.Join(constants.DefaultConfigDir, "edgeagent.yaml"),
	}
}